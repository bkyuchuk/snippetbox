package mocks

import "snippetbox.bogdandev.de/internal/models"

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "doe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (int64, error) {
	if email == "bogdan@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int64) (bool, error) {
	if id == 1 {
		return true, nil
	}

	return false, nil
}
