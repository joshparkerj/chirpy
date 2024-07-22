package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func validateChirp(res http.ResponseWriter, req *http.Request) {
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

	unprofaneRegexp := regexp.MustCompile(`(?i)(^|\s)(kerfuffle|sharbert|fornax)($|\s)`)
	unprofane := unprofaneRegexp.ReplaceAllString(params.Body, "$1****$3")
	unprofane = unprofaneRegexp.ReplaceAllString(unprofane, "$1****$3")

	respBody := cleanedResponse{
		CleanedBody: unprofane,
	}

	sendJsonResponse(respBody, res)
}
