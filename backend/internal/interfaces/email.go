package interfaces

type EmailService interface {
	ConfirmEmail(email, token string) error
}
