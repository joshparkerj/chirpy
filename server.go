package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

func serverPortHelper(server *http.Server, changePort func(oldPort int) (newPort int)) {
	port, err := strconv.Atoi(server.Addr[1:])
	if err != nil {
		fmt.Println(err)
		port = 8080
	} else {
		port = changePort(port)
	}

	server.Addr = fmt.Sprintf(":%d", port)
}

func jiggerPort(server *http.Server) {
	serverPortHelper(server, func(oldPort int) (newPort int) {
		newPort = oldPort
		return
	})
}

func rejiggerPort(server *http.Server) {
	serverPortHelper(server, func(oldPort int) (newPort int) {
		roll := rand.Intn(6) + 1
		newPort = oldPort + roll
		return
	})
}
