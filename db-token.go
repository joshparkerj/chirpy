package main

import (
	"errors"
	"fmt"
	"time"
)

func (db *DB) CreateToken(userId int) (token string, err error) {
	token, err = generateRefreshToken()
	if err != nil {
		return
	}

	err = db.ensureDB()
	if err != nil {
		return
	}

	dur, err := time.ParseDuration(fmt.Sprintf("%dh", 24*60))
	if err != nil {
		return
	}

	expiry := time.Now().Add(dur)

	tokenRecord := Token{
		Expiry:       expiry,
		RefreshToken: token,
		UserID:       userId,
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return
	}

	tokenRecord.ID = len(databaseStructure.Tokens) + 1

	databaseStructure.Tokens[tokenRecord.ID] = tokenRecord
	db.writeDB(databaseStructure)
	return
}

func (db *DB) GetToken(token string) (matchingToken Token, err error) {
	err = db.ensureDB()
	if err != nil {
		return
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return
	}

	for _, matchingToken = range databaseStructure.Tokens {
		if matchingToken.RefreshToken == token {
			return
		}
	}

	err = errors.New("token not found")
	return
}

func (db *DB) UpdateToken(token Token) (err error) {
	err = db.ensureDB()
	if err != nil {
		return
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return
	}

	databaseStructure.Tokens[token.ID] = token

	db.writeDB(databaseStructure)
	return
}
