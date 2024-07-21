package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const adminHtmlTemplate = `
			<html><body>
			<h1>Welcome, Chirpy Admin</h1>
			%s
			</body></html>
`

type errorResponse struct {
	Error string `json:"error"`
}

type apiConfig struct {
	fileserverHits int
}

type metricsHandler struct {
	nextHandler http.Handler
	apicfg      *apiConfig
}

func (mh *metricsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	mh.apicfg.fileserverHits++
	fmt.Printf("counted hit number: %d\n", mh.apicfg.fileserverHits)
	mh.nextHandler.ServeHTTP(res, req)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return &metricsHandler{
		nextHandler: next,
		apicfg:      cfg,
	}
}

func (cfg *apiConfig) getMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) htmlMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(200)
	message := fmt.Sprintf("<p>Chirpy has been visited %d times!</p>", cfg.fileserverHits)
	res.Write([]byte(fmt.Sprintf(adminHtmlTemplate, message)))
}

func (cfg *apiConfig) resetMetrics(res http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = 0
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("reset"))
}

func somethingWentWrong(res http.ResponseWriter) {
	sendErrorResponse("Something went wrong", 500, res)
}

func sendErrorResponse(err string, statusCode int, res http.ResponseWriter) {
	respBody := errorResponse{
		Error: err,
	}

	dat, _ := json.Marshal(respBody)
	res.WriteHeader(statusCode)
	res.Write(dat)
}

func main() {
	apicfg := apiConfig{fileserverHits: 0}
	sm := http.NewServeMux()

	sm.HandleFunc("GET /api/healthz", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})

	sm.Handle("/app/*", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	sm.HandleFunc("GET /api/metrics", apicfg.getMetrics)

	sm.HandleFunc("GET /admin/metrics", apicfg.htmlMetrics)

	sm.HandleFunc("/api/reset", apicfg.resetMetrics)

	sm.HandleFunc("POST /api/validate_chirp", func(res http.ResponseWriter, req *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}

		type validResponse struct {
			Valid bool `json:"valid"`
		}

		decoder := json.NewDecoder(req.Body)
		params := parameters{}
		err := decoder.Decode(&params)

		if err != nil {
			somethingWentWrong(res)
			return
		}

		if len(params.Body) > 140 {
			sendErrorResponse("Chirp is too long", 400, res)
			return
		}

		respBody := validResponse{
			Valid: true,
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			somethingWentWrong(res)
			return
		}

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(dat)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	server.ListenAndServe()
}
