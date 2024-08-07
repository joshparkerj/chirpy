package main

import (
	"crypto/rand"
	"encoding/hex"
)

func generateRefreshToken() (token string, err error) {
	data := make([]byte, 32)
	_, err = rand.Read(data)
	if err != nil {
		return
	}

	token = hex.EncodeToString(data)
	return
}
