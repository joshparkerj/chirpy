package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleApiError(err error, msg string, statusCode int, res http.ResponseWriter) {
	log.Default().Println(msg)
	log.Default().Println(err)
	errorString := fmt.Sprintf("%s %s", msg, err)
	sendErrorResponse(errorString, statusCode, res)
}

func consoleLog(msg string) {
	fmt.Println(msg)
}
