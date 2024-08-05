package main

import (
	"encoding/json"
	"net/http"
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

	sendOkJsonResponse(chirps, res)
}

func postChirp(res http.ResponseWriter, req *http.Request) {
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

	chirp, err := db.CreateChirp(unprofane)
	if err != nil {
		handleApiError(err, "error in CreateChirp", 500, res)
		return
	}

	sendJsonResponse(chirp, res, 201)
}
