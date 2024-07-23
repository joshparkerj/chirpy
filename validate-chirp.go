package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

func cleanParams(params parameters, res http.ResponseWriter) (string, error) {
	if len(params.Body) > 140 {
		sendErrorResponse("Chirp is too long", 400, res)
		return "", errors.New("body too long")
	}

	unprofaneRegexp := regexp.MustCompile(`(?i)(^|\s)(kerfuffle|sharbert|fornax)($|\s)`)
	unprofane := unprofaneRegexp.ReplaceAllString(params.Body, "$1****$3")
	unprofane = unprofaneRegexp.ReplaceAllString(unprofane, "$1****$3")
	return unprofane, nil
}

func validateChirp(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		somethingWentWrong(res)
		return
	}

	unprofane, err := cleanParams(params, res)
	if err != nil {
		return
	}

	respBody := cleanedResponse{
		CleanedBody: unprofane,
	}

	sendOkJsonResponse(respBody, res)
}
