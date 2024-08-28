package mocks

import (
	"github.com/sotnikea/Go_Learn/tree/main/snippetbox/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) Authenticate(email, password string) (interface{}, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return "111111111111111111111111", nil
	}
	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id string) (bool, error) {
	switch id {
	case "111111111111111111111111":
		return true, nil
	default:
		return false, nil
	}
}
