package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"sync"
)

const dbFilename = "database.json"

type DB struct {
	path string
	mux  *sync.Mutex
}

// newDB creates a database connection
// and creates the database file if it doesn't exist
func newDB(path string) (*DB, error) {
	db := DB{
		path,
		&sync.Mutex{},
	}

	db.mux.Lock()

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disc
func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirp := Chirp{}
	chirp.Body = body
	err := db.ensureDB()
	if err != nil {
		return chirp, err
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	chirp.ID = len(databaseStructure.Chirps) + 1
	databaseStructure.Chirps[chirp.ID] = chirp
	db.writeDB(databaseStructure)
	return chirp, nil
}

func find(chirps []Chirp, id int) *Chirp {
	for _, chirp := range chirps {
		if chirp.ID == id {
			return &chirp
		}
	}

	return nil
}

// GetChirps returns all the chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	err := db.ensureDB()
	if err != nil {
		return make([]Chirp, 0), err
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return make([]Chirp, 0), err
	}

	chirps := make([]Chirp, len(databaseStructure.Chirps))
	for i, chirp := range databaseStructure.Chirps {
		chirps[i-1] = chirp
	}

	return chirps, nil
}

func createNewDatabaseFile() error {
	dbStructure := DBStructure{}
	dbStructure.Chirps = make(map[int]Chirp, 0)
	newContent, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(dbFilename, newContent, 0666)
	if err != nil {
		return err
	}

	return nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(dbFilename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			err := createNewDatabaseFile()
			if err != nil {
				log.Default().Println("error in create new database")
				return err
			}

			return nil
		}

		return err
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	stubStructure := DBStructure{}
	fileReader, err := os.Open(dbFilename)
	if err != nil {
		return stubStructure, err
	}

	decoder := json.NewDecoder(fileReader)
	err = decoder.Decode(&stubStructure)
	if err != nil {
		return stubStructure, err
	}

	return stubStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	databaseContent, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(dbFilename, databaseContent, 0666)
	if err != nil {
		return err
	}

	return nil
}
