package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

func createJwt(expirySeconds, userId int, secret string) (string, error) {
	method := jwt.SigningMethodHS256
	audience := jwt.ClaimStrings(make([]string, 0))
	startTime := time.Now()
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", expirySeconds))
	if err != nil {
		return "", err
	}

	expiry := startTime.Add(dur)
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   fmt.Sprintf("%d", userId),
		Audience:  audience,
		ExpiresAt: jwt.NewNumericDate(expiry),
		NotBefore: jwt.NewNumericDate(startTime),
		IssuedAt:  jwt.NewNumericDate(startTime),
	}

	token := jwt.NewWithClaims(method, claims)
	return token.SignedString([]byte(secret))
}

func validateJwt(token, key string) (userId int, err error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(key), nil })
	if err != nil {
		return
	}

	claims := parsedToken.Claims
	userIdString, err := claims.GetSubject()
	if err != nil {
		return
	}

	userId, err = strconv.Atoi(userIdString)
	return
}
