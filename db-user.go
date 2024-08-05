package main

func (db *DB) CreateUser(email string, password string) (user User, err error) {
	user = User{}
	user.Email = email
	user.Password = password
	err = db.ensureDB()
	if err != nil {
		return
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return
	}

	user.ID = len(databaseStructure.Users) + 1
	databaseStructure.Users[user.ID] = user
	db.writeDB(databaseStructure)
	return
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

func (db *DB) UpdateUser(email, password string, userId int) (user User, err error) {
	user = User{
		Email:    email,
		Password: password,
		ID:       userId,
	}

	err = db.ensureDB()
	if err != nil {
		return
	}

	databaseStructure, err := db.loadDB()
	if err != nil {
		return
	}

	databaseStructure.Users[user.ID] = user
	db.writeDB(databaseStructure)
	return
}
