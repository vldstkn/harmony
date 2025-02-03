package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"harmony/internal/config"
)

type Service struct {
	Config *config.Config
}
type ServiceDeps struct {
	Config *config.Config
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Config: deps.Config,
	}
}

func (service *Service) ConfirmEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "harmony@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Подтверждение почты")
	m.SetBody("text/html", fmt.Sprintf("Здравствуйте! <b>Подтвердите вашу почту</b>, "+
		"перейдя по ссылке: <a href='https://flotnend.com/confirm?token=%s'>Подтвердить</a><br> <strong>%s</strong>", token, token))
	d := gomail.NewDialer("smtp.gmail.com", 587, service.Config.EmailLogin, service.Config.EmailPass)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
