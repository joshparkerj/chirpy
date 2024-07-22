package main

import (
	"encoding/json"
	"net/http"
)

func somethingWentWrong(res http.ResponseWriter) {
	sendErrorResponse("Something went wrong", 500, res)
}

func plainTextResponse(response string, res http.ResponseWriter) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	res.Write([]byte(response))
}

func sendErrorResponse(err string, statusCode int, res http.ResponseWriter) {
	respBody := errorResponse{
		Error: err,
	}

	dat, _ := json.Marshal(respBody)
	res.WriteHeader(statusCode)
	res.Write(dat)
}

func sendJsonResponse(respBody interface{}, res http.ResponseWriter) {
	dat, err := json.Marshal(respBody)
	if err != nil {
		somethingWentWrong(res)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(dat)
}
