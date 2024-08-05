package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return &metricsHandler{
		nextHandler: next,
		apicfg:      cfg,
	}
}

func (cfg *apiConfig) getMetrics(res http.ResponseWriter, req *http.Request) {
	plainTextResponse(fmt.Sprintf("Hits: %d", cfg.fileserverHits), res)
}

func (cfg *apiConfig) htmlMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(200)
	message := fmt.Sprintf("<p>Chirpy has been visited %d times!</p>", cfg.fileserverHits)
	res.Write([]byte(fmt.Sprintf(adminHtmlTemplate, message)))
}

func (cfg *apiConfig) resetMetrics(res http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = 0
	plainTextResponse("reset", res)
}
