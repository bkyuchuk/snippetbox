package models

import (
	"database/sql"
	"time"
)

type User struct {
	Id             int64
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	now := time.Now().UTC()

	_, err := m.DB.Exec(`INSERT INTO users (name, email, hashed_password, created)
		VALUES (?, ?, ?, ?)
	`, name, email, password, now)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int64, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int64) (bool, error) {
	return false, nil
}
