package api

import (
	"net/http"
)

type Service struct {
	JWTSecret string
}

func NewService(JWTSecret string) *Service {
	return &Service{
		JWTSecret: JWTSecret,
	}
}

func (service *Service) AddCookie(w *http.ResponseWriter, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   maxAge,
	}
	http.SetCookie(*w, cookie)
}
