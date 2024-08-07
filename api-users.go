package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func prepUser(res http.ResponseWriter, req *http.Request) (db *DB, hashedPassword []byte, reqNewUser newUser, err error) {
	db, err = newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqNewUser = newUser{}
	err = decoder.Decode(&reqNewUser)

	if err != nil {
		handleApiError(err, "error in Decode", 500, res)
		return
	}

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(reqNewUser.Password), 10)
	if err != nil {
		handleApiError(err, "error in GenerateFromPassword", 500, res)
		return
	}

	return
}

func postUser(res http.ResponseWriter, req *http.Request) {
	db, hashedPassword, reqNewUser, err := prepUser(res, req)
	if err != nil {
		return
	}

	user, err := db.CreateUser(reqNewUser.Email, string(hashedPassword))
	if err != nil {
		handleApiError(err, "error in CreateUser", 500, res)
		return
	}

	resUser := userPasswordRedacted{
		Email: user.Email,
		ID:    user.ID,
	}

	sendJsonResponse(resUser, res, 201)
}

func putUser(res http.ResponseWriter, req *http.Request) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) < 7 {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	authorization := authHeader[7:]
	userId, err := validateJwt(authorization, apicfg.jwtSecret)
	if err != nil {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	db, hashedPassword, reqNewUser, err := prepUser(res, req)
	if err != nil {
		return
	}

	user, err := db.UpdateUser(reqNewUser.Email, string(hashedPassword), userId)

	if err != nil {
		handleApiError(err, "error in UpdateUser", 500, res)
		return
	}

	resUser := userPasswordRedacted{
		Email: user.Email,
		ID:    user.ID,
	}

	sendJsonResponse(resUser, res, 200)
}

func postLogin(res http.ResponseWriter, req *http.Request) {
	db, err := newDB(dbFilename)
	if err != nil {
		handleApiError(err, "error in newDB", 500, res)
	}

	decoder := json.NewDecoder(req.Body)
	reqUser := newUser{}
	err = decoder.Decode(&reqUser)

	if err != nil {
		handleApiError(err, "error in Decode", 500, res)
		return
	}

	dbUser, err := db.GetUser(reqUser.Email)
	if err != nil {
		handleApiError(err, "error in GetUser", 500, res)
		return
	}

	if dbUser == nil {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(reqUser.Password))
	if err != nil {
		handleApiError(nil, "unauthorized", 401, res)
		return
	}

	oneHour := 60 * 60
	var expiry int
	if reqUser.Expiry > 0 && reqUser.Expiry <= oneHour {
		expiry = reqUser.Expiry
	} else {
		expiry = oneHour
	}

	token, err := createJwt(expiry, dbUser.ID, apicfg.jwtSecret)
	if err != nil {
		handleApiError(err, "error in createJwt", 500, res)
		return
	}

	refreshToken, err := db.CreateToken(dbUser.ID)
	if err != nil {
		handleApiError(err, "error in CreateToken", 500, res)
		return
	}

	resUser := userPasswordRedacted{
		Email:        dbUser.Email,
		ID:           dbUser.ID,
		Token:        token,
		RefreshToken: refreshToken,
	}

	sendJsonResponse(resUser, res, 200)
}
