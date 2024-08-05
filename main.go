package main

import (
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

var apicfg apiConfig

func main() {
	godotenv.Load()
	apicfg = apiConfig{
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	sm := http.NewServeMux()

	sm.HandleFunc("GET /api/healthz", func(res http.ResponseWriter, req *http.Request) {
		plainTextResponse("OK", res)
	})

	sm.Handle("/app/*", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	sm.HandleFunc("GET /api/metrics", apicfg.getMetrics)

	sm.HandleFunc("GET /admin/metrics", apicfg.htmlMetrics)

	sm.HandleFunc("/api/reset", apicfg.resetMetrics)

	sm.HandleFunc("POST /api/validate_chirp", validateChirp)

	sm.HandleFunc("GET /api/chirps/{chirpID}", getChirp)

	sm.HandleFunc("GET /api/chirps", getChirps)

	sm.HandleFunc("POST /api/chirps", postChirp)

	sm.HandleFunc("POST /api/users", postUser)

	sm.HandleFunc("PUT /api/users", putUser)

	sm.HandleFunc("POST /api/login", postLogin)

	sm.HandleFunc("/", handleHomePage)

	defaultPort := os.Getenv("PORT")
	server := &http.Server{
		Addr:    ":" + defaultPort,
		Handler: sm,
	}

	jiggerPort(server)

	consoleLog("now about to listen and serve on " + server.Addr)
	for err := server.ListenAndServe(); err != nil; {
		consoleLog(err.Error())
		rejiggerPort(server)
		consoleLog("now about to listen and serve on " + server.Addr)
		err = server.ListenAndServe()
	}
}
