package main

import (
	"errors"
	"net/http"
)

func authorize(req *http.Request) (userId int, err error) {
	authHeader := req.Header.Get("Authorization")
	if len(authHeader) < 7 {
		err = errors.New("no authorization in header")
		return
	}

	authorization := authHeader[7:]
	userId, err = validateJwt(authorization, apicfg.jwtSecret)
	if err != nil {
		return
	}

	return
}
