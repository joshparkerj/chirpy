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

func createNewDatabaseFile() error {
	dbStructure := DBStructure{}

	// don't forget to initialize database structures here!
	// TODO: find a way to eliminate the need to initialize each database structure here
	// (because it's too easy to forget when adding addtional database structures)
	dbStructure.Chirps = make(map[int]Chirp, 0)
	dbStructure.Users = make(map[int]User, 0)
	dbStructure.Tokens = make(map[int]Token, 0)

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
