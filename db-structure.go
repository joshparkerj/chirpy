package main

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
	Tokens map[int]Token `json:"tokens"`
}
