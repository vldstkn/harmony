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

func (service *Service) Register(email, password, name string) (string, error) {
	err := service.Repository.UserNotExists(email, name)
	if err != nil {
		if err.Error() == "name" {
			return "", errors.New("the name is already taken")
		} else {
			return "", errors.New("there is already an account with this email")
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
		return "", errors.New("internal error")
	}
	return id, nil
}
func (service *Service) Login(email, password string) (string, error) {
	user := service.Repository.GetByEmail(email)
	if user == nil {
		return "", errors.New("invalid email or password")
	}
	if user.Status == models.Unconfirmed {
		// TODO:
		return "", errors.New("confirm your email address")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	return user.Id, nil
}

func (service *Service) ConfirmEmail(secret, token string) (string, error) {
	isValid, data := jwt.NewJWT(secret).Parse(token)
	if !isValid {
		service.Logger.Error("jwt.NewJWT(secret).Parse", slog.String("err", "token not valid"))
		return "", errors.New("token not valid")
	}
	err := service.Repository.UpdateById("status", string(models.Confirmed), data.Id)
	if err != nil {
		service.Logger.Error("Repository.UpdateById", slog.String("err", err.Error()))
		return "", errors.New("internal error")
	}
	return data.Id, nil
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
