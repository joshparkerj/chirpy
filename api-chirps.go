package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
)

func getChirp(res http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	db, err := newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		handleApiError(err, "error in GetChirps", 500, res)
		return
	}

	idNum, err := strconv.Atoi(chirpID)
	if err != nil {
		handleApiError(err, "error in Atoi", 500, res)
		return
	}

	chirp := find(chirps, idNum)
	if chirp == nil {
		res.WriteHeader(404)
	} else {
		sendOkJsonResponse(chirp, res)
	}
}

func getChirps(res http.ResponseWriter, req *http.Request) {
	db, err := newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		handleApiError(err, "error in GetChirps", 500, res)
		return
	}

	var cmp func(a, b Chirp) int
	sort := req.URL.Query().Get("sort")
	if sort == "desc" {
		cmp = func(a, b Chirp) int {
			return b.ID - a.ID
		}
	} else {
		cmp = func(a, b Chirp) int {
			return a.ID - b.ID
		}
	}

	authorId := req.URL.Query().Get("author_id")
	if authorId != "" {
		idNum, err := strconv.Atoi(authorId)
		if err != nil {
			handleApiError(err, "error in Atoi", 500, res)
			return
		}

		filteredChirps := make([]Chirp, 0)
		for _, chirp := range chirps {
			if chirp.AuthorId == idNum {
				filteredChirps = append(filteredChirps, chirp)
			}
		}

		slices.SortFunc(filteredChirps, cmp)
		sendOkJsonResponse(filteredChirps, res)
	} else {
		slices.SortFunc(chirps, cmp)
		sendOkJsonResponse(chirps, res)
	}
}

func postChirp(res http.ResponseWriter, req *http.Request) {
	userId, err := authorize(req)
	if err != nil {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	db, err := newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		handleApiError(err, "error in Decode", 500, res)
		return
	}

	unprofane, err := cleanParams(params, res)
	if err != nil {
		return
	}

	chirp, err := db.CreateChirp(unprofane, userId)
	if err != nil {
		handleApiError(err, "error in CreateChirp", 500, res)
		return
	}

	sendJsonResponse(chirp, res, 201)
}

func deleteChirp(res http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	userId, err := authorize(req)
	if err != nil {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	db, err := newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
		return
	}

	idNum, err := strconv.Atoi(chirpID)
	if err != nil {
		handleApiError(err, "error in Atoi", 500, res)
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		handleApiError(err, "error in GetChirps", 500, res)
		return
	}

	chirp := find(chirps, idNum)
	if chirp == nil {
		handleApiError(nil, "not found", 404, res)
		return
	} else if chirp.AuthorId != userId {
		handleApiError(nil, "not authorized", 403, res)
		return
	}

	err = db.DeleteChirp(idNum)
	if err != nil {
		handleApiError(err, "error in DeleteChirp", 500, res)
		return
	}

	sendJsonResponse(nil, res, 204)
}
