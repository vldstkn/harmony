package interfaces

import (
	"harmony/internal/models"
	"harmony/pkg/jwt"
)

type AccountService interface {
	Register(email, password, name string) (string, error)
	Login(email, password string) (string, error)
	IssueTokens(secret string, data jwt.Data) (string, string, error)
	ConfirmEmail(secret, token string) (string, error)
}
type AccountRepository interface {
	GetById(id string) *models.User
	GetByEmail(id string) *models.User
	Create(user models.User) (string, error)

	DeleteUnique(field, value string) error
	SearchByName(name string) []models.User
	UpdateById(field, value, id string) error
	UserNotExists(email, name string) error
}
