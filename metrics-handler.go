package main

import (
	"fmt"
	"net/http"
)

type metricsHandler struct {
	nextHandler http.Handler
	apicfg      *apiConfig
}

func (mh *metricsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	mh.apicfg.fileserverHits++
	fmt.Printf("counted hit number: %d\n", mh.apicfg.fileserverHits)
	mh.nextHandler.ServeHTTP(res, req)
}
