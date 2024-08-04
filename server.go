package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

func rejiggerPort(server *http.Server) (err error) {
	port, err := strconv.Atoi(server.Addr[1:])
	if err != nil {
		return
	}
	roll := rand.Intn(6) + 1
	port += roll
	server.Addr = fmt.Sprintf(":%d", port)
	return
}
