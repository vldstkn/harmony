package interfaces

import (
	"harmony/internal/models"
	"harmony/pkg/jwt"
)

type AccountService interface {
	Register(email, password, name string) (int64, error)
	Login(email, password string) (int64, error)
	IssueTokens(secret string, data jwt.Data) (string, string, error)
	ConfirmEmail(secret, token string) (int64, error)

	FindByName(id int64, name string) []models.User
	AddFriend(userId, friendId int64) error
	DeleteFriend(userId, friendId int64) error
	FindFriendsByName(userId int64, name string) []models.User
}

type AccountRepository interface {
	GetById(id int64) *models.User
	GetByEmail(email string) *models.User
	Create(user models.User) (int64, error)
	DeleteById(id int64) error
	SearchByName(name string) []models.User
	UpdateStatusById(value string, id int64) error
	UserNotExists(email, name string) error

	FindByName(id int64, name string) []models.User
	AddFriend(userId, friendId int64) error
	DeleteFriend(userId, friendId int64) error
	FindFriendsByName(userId int64, name string) []models.User
}
