package main

import (
	"net/http"
	"fmt"
)

type apiConfig struct {
	fileserverHits int
}

type metricsHandler struct {
	nextHandler http.Handler
	apicfg *apiConfig
}

func (mh *metricsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	mh.apicfg.fileserverHits++
	fmt.Printf("counted hit number: %d", mh.apicfg.fileserverHits)
	mh.nextHandler.ServeHTTP(res, req)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return &metricsHandler{
		nextHandler: next,
		apicfg: cfg,
	}
}

func (cfg *apiConfig) getMetrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) resetMetrics(res http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = 0
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte("reset"))
}

func main() {
	apicfg := apiConfig{ fileserverHits: 0 }
	sm := http.NewServeMux()

	sm.HandleFunc("/healthz", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(200)
		res.Write([]byte("OK"))
	})

	sm.Handle("/app/*", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	sm.HandleFunc("/metrics", apicfg.getMetrics)

	sm.HandleFunc("/reset", apicfg.resetMetrics)

        server := &http.Server{
		Addr: ":8080",
		Handler: sm,
	}

	server.ListenAndServe()
}
