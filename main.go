package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
			log.Default().Println("error in newDB")
			log.Default().Println(err)
			somethingWentWrong(res)
			return
		}

		chirps, err := db.GetChirps()
		if err != nil {
			log.Default().Println("error in GetChirps")
			log.Default().Println(err)
			somethingWentWrong(res)
			return
		}

		idNum, err := strconv.Atoi(chirpID)
		if err != nil {
			log.Default().Println("error in Atoi")
			log.Default().Println(err)
			somethingWentWrong(res)
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
			log.Default().Println("error in newDB")
			log.Default().Println(err)
			somethingWentWrong(res)
			return
		}

		chirps, err := db.GetChirps()
		if err != nil {
			log.Default().Println("error in GetChirps")
			log.Default().Println(err)
			somethingWentWrong(res)
			return
		}

		sendOkJsonResponse(chirps, res)
	})

	sm.HandleFunc("POST /api/chirps", func(res http.ResponseWriter, req *http.Request) {
		db, err := newDB(dbFilename)
		if err != nil {
			log.Default().Println("error in newDB")
			log.Default().Println(err)
			somethingWentWrong(res)
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err = decoder.Decode(&params)

		if err != nil {
			log.Default().Println("error in Decode")
			log.Default().Println(err)
			somethingWentWrong(res)
			return
		}

		unprofane, err := cleanParams(params, res)
		if err != nil {
			return
		}

		chirp, err := db.CreateChirp(unprofane)
		if err != nil {
			log.Default().Println("error in CreateChirp")
			log.Default().Println(err)
			somethingWentWrong(res)
			return
		}

		sendJsonResponse(chirp, res, 201)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	server.ListenAndServe()
}
