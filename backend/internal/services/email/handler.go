package email

import (
	"context"
	"harmony/internal/config"
	"harmony/internal/interfaces"
	pb "harmony/pkg/api/email"
	"log/slog"
)

type Handler struct {
	pb.UnsafeEmailServer
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.EmailService
}
type HandlerDeps struct {
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.EmailService
}

func NewHandler(deps HandlerDeps) *Handler {
	return &Handler{
		Config:  deps.Config,
		Logger:  deps.Logger,
		Service: deps.Service,
	}
}

func (handler *Handler) ConfirmEmail(ctx context.Context, r *pb.ConfirmEmailSendReq) (*pb.ConfirmEmailSendRes, error) {
	err := handler.Service.ConfirmEmail(r.Email, r.Token)
	return nil, err
}
