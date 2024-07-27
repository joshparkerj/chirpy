package main

func (db *DB) CreateUser(email string, password string) (User, error) {
	user := User{}
	user.Email = email
	user.Password = password
	err := db.ensureDB()
	if err != nil {
		return user, err
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return user, err
	}

	user.ID = len(databaseStructure.Users) + 1
	databaseStructure.Users[user.ID] = user
	db.writeDB(databaseStructure)
	return user, nil
}

func (db *DB) GetUser(email string) (*User, error) {
	err := db.ensureDB()
	if err != nil {
		return &User{}, err
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return &User{}, err
	}

	for _, user := range databaseStructure.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, nil
}
