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

type newUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Expiry   int    `json:"expires_in_seconds"`
}

type userPasswordRedacted struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
