package main

import (
	"net/http"
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	server.ListenAndServe()
}
