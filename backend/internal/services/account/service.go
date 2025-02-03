package account

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"harmony/internal/interfaces"
	"harmony/internal/models"
	"harmony/pkg/jwt"
	"log/slog"
	"time"
)

type ServiceDeps struct {
	Repository interfaces.AccountRepository
	Logger     *slog.Logger
}
type Service struct {
	Repository interfaces.AccountRepository
	Logger     *slog.Logger
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Repository: deps.Repository,
		Logger:     deps.Logger,
	}
}

func (service *Service) Register(email, password, name string) (int64, error) {
	err := service.Repository.UserNotExists(email, name)
	if err != nil {
		if err.Error() == "name" {
			return -1, errors.New("the name is already taken")
		} else {
			return -1, errors.New("there is already an account with this email")
		}
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Email:    email,
		Name:     name,
		Password: string(hashPassword),
	}
	id, err := service.Repository.Create(user)
	if err != nil {
		return -1, errors.New("internal error")
	}
	return id, nil
}
func (service *Service) Login(email, password string) (int64, error) {
	user := service.Repository.GetByEmail(email)
	if user == nil {
		return -1, errors.New("invalid email or password")
	}
	if user.Status == models.Unconfirmed {
		// TODO:
		return -1, errors.New("confirm your email address")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return -1, errors.New("invalid email or password")
	}
	return user.Id, nil
}

func (service *Service) ConfirmEmail(secret, token string) (int64, error) {
	isValid, data := jwt.NewJWT(secret).Parse(token)
	if !isValid {
		service.Logger.Error("jwt.NewJWT(secret).Parse", slog.String("err", "token not valid"))
		return -1, errors.New("token not valid")
	}
	err := service.Repository.UpdateStatusById(string(models.Confirmed), data.Id)
	if err != nil {
		service.Logger.Error("Repository.UpdateById", slog.String("err", err.Error()))
		return -1, errors.New("internal error")
	}
	return data.Id, nil
}

func (service *Service) FindByName(id int64, name string) []models.User {
	users := service.Repository.FindByName(id, name)
	return users
}
func (service *Service) AddFriend(userId, friendId int64) error {
	if userId == friendId {
		return errors.New("bad request")
	}
	err := service.Repository.AddFriend(userId, friendId)
	if err != nil {
		service.Logger.Error("Repository.AddFriend", slog.String("err", err.Error()))
		return errors.New("bad request")
	}
	return nil
}
func (service *Service) DeleteFriend(userId, friendId int64) error {
	err := service.Repository.DeleteFriend(userId, friendId)
	if err != nil {
		return errors.New("bad request")
	}
	return nil
}
func (service *Service) FindFriendsByName(userId int64, name string) []models.User {
	friends := service.Repository.FindFriendsByName(userId, name)
	return friends
}

func (service *Service) IssueTokens(secret string, data jwt.Data) (string, string, error) {
	j := jwt.NewJWT(secret)
	accessToken, err := j.Create(data, time.Now().Add(time.Hour*2).Add(time.Minute*10))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := j.Create(data, time.Now().AddDate(0, 0, 2).Add(time.Hour*2))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
