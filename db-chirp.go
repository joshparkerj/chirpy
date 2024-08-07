package main

// CreateChirp creates a new chirp and saves it to disc
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
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
	chirp.AuthorId = authorId
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

func (db *DB) DeleteChirp(chirpId int) (err error) {
	err = db.ensureDB()
	if err != nil {
		return
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return
	}

	delete(databaseStructure.Chirps, chirpId)
	db.writeDB(databaseStructure)
	return
}
