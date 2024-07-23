#! /bin/bash
rm -f database.json && gofmt -w . && staticcheck && gotags ./*.go > tags && go run .
