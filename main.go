package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	apicfg := apiConfig{fileserverHits: 0}
	sm := http.NewServeMux()

	sm.HandleFunc("GET /api/healthz", func(res http.ResponseWriter, req *http.Request) {
		plainTextResponse("OK", res)
	})

	sm.Handle("/app/*", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	sm.HandleFunc("GET /api/metrics", apicfg.getMetrics)

	sm.HandleFunc("GET /admin/metrics", apicfg.htmlMetrics)

	sm.HandleFunc("/api/reset", apicfg.resetMetrics)

	sm.HandleFunc("POST /api/validate_chirp", validateChirp)

	sm.HandleFunc("GET /api/chirps/{chirpID}", func(res http.ResponseWriter, req *http.Request) {
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
	})

	sm.HandleFunc("GET /api/chirps", func(res http.ResponseWriter, req *http.Request) {
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
	})

	sm.HandleFunc("POST /api/chirps", func(res http.ResponseWriter, req *http.Request) {
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
	})

	sm.HandleFunc("POST /api/users", func(res http.ResponseWriter, req *http.Request) {
		db, err := newDB(dbFilename)
		if err != nil {
			handleApiError(err, "error in newDB", 500, res)
			return
		}

		decoder := json.NewDecoder(req.Body)
		reqNewUser := newUser{}
		err = decoder.Decode(&reqNewUser)

		if err != nil {
			handleApiError(err, "error in Decode", 500, res)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqNewUser.Password), 10)
		if err != nil {
			handleApiError(err, "error in GenerateFromPassword", 500, res)
			return
		}

		user, err := db.CreateUser(reqNewUser.Email, string(hashedPassword))
		if err != nil {
			handleApiError(err, "error in CreateUser", 500, res)
			return
		}

		resUser := userPasswordRedacted{
			Email: user.Email,
			ID:    user.ID,
		}

		sendJsonResponse(resUser, res, 201)
	})

	sm.HandleFunc("POST /api/login", func(res http.ResponseWriter, req *http.Request) {
		db, err := newDB(dbFilename)
		if err != nil {
			handleApiError(err, "error in newDB", 500, res)
		}

		decoder := json.NewDecoder(req.Body)
		reqUser := newUser{}
		err = decoder.Decode(&reqUser)

		if err != nil {
			handleApiError(err, "error in Decode", 500, res)
			return
		}

		dbUser, err := db.GetUser(reqUser.Email)
		if err != nil {
			handleApiError(err, "error in GetUser", 500, res)
			return
		}

		if dbUser == nil {
			handleApiError(nil, "unauthorized", 401, res)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(reqUser.Password))
		if err != nil {
			handleApiError(nil, "unauthorized", 401, res)
			return
		}

		resUser := userPasswordRedacted{
			Email: dbUser.Email,
			ID:    dbUser.ID,
		}

		sendJsonResponse(resUser, res, 200)
	})

	sm.HandleFunc("/", handleHomePage)

	server := &http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	consoleLog("now about to listen and serve on " + server.Addr)
	for err := server.ListenAndServe(); err != nil; {
		consoleLog(err.Error())
		rejiggerPort(server)
		consoleLog("now about to listen and serve on " + server.Addr)
		err = server.ListenAndServe()
	}
}
