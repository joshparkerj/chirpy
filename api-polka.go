package main

import (
	"encoding/json"
	"net/http"
)

func polkaWebhooks(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) < 7 {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	authorization := authHeader[7:]
	if authorization != apicfg.polkaApiKey {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	decoder := json.NewDecoder(req.Body)
	whEvent := polkaEvent{}
	err := decoder.Decode(&whEvent)
	if err != nil {
		handleApiError(err, "error in Decode", 500, res)
		return
	}

	if whEvent.Event == "user.upgraded" {
		userId := whEvent.Data.UserID
		db, err := newDB(dbFilename)
		if err != nil {
			handleApiError(err, "error in newDB", 500, res)
			return
		}

		err = db.UpgradeUser(userId)
		if err != nil {
			if err.Error() == "user not found" {
				handleApiError(err, "could not upgrade. user id does not exist.", 404, res)
				return
			} else {
				handleApiError(err, "error in UpgradeUser", 500, res)
				return
			}
		}

		sendJsonResponse(nil, res, 204)
	} else {
		sendJsonResponse(nil, res, 204)
	}
}
