package dto

type AccountLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AccountLoginRes struct {
	Id          string `json:"id"`
	AccessToken string `json:"access_token"`
}

type AccountRegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type AccountConfirmEmailReq struct {
	Token string `json:"token" validate:"required"`
}
type AccountConfirmEmailRes struct {
	Id          string `json:"id"`
	AccessToken string `json:"access_token"`
}
