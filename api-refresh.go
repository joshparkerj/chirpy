package main

import (
	"fmt"
	"net/http"
	"time"
)

func postRefresh(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) < 7 {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	authorization := authHeader[7:]
	db, err := newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
		return
	}

	token, err := db.GetToken(authorization)
	if err != nil {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	fmt.Println(time.Until(token.Expiry))

	if token.Expiry.Before(time.Now()) {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	user, err := db.GetUserByID(token.UserID)

	if err != nil {
		handleApiError(err, "err in GetUserByID", 500, res)
		return
	}

	jwt, err := createJwt(60*60, user.ID, apicfg.jwtSecret)
	if err != nil {
		handleApiError(err, "error in createJwt", 500, res)
		return
	}

	resToken := tokenResponse{
		Token: jwt,
	}

	sendJsonResponse(resToken, res, 200)
}

func revoke(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) < 7 {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	authorization := authHeader[7:]
	db, err := newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
		return
	}

	token, err := db.GetToken(authorization)
	if err != nil || token.Expiry.Before(time.Now()) {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	token.Expiry = time.Now()
	db.UpdateToken(token)
	sendJsonResponse(nil, res, 204)
}
