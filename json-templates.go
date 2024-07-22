package main

type errorResponse struct {
	Error string `json:"error"`
}

type parameters struct {
	Body string `json:"body"`
}

type cleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}
