package main

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}
