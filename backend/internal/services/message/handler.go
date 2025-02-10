package messages

import (
	"context"
	"harmony/internal/config"
	"harmony/internal/interfaces"
	pbmes "harmony/pkg/api/messages"
	"log/slog"
)

type Handler struct {
	pbmes.UnsafeMessagesServer
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.MessageService
}

type HandlerDeps struct {
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.MessageService
}

func NewHandler(deps *HandlerDeps) *Handler {
	return &Handler{
		Config:  deps.Config,
		Logger:  deps.Logger,
		Service: deps.Service,
	}
}
func (handler *Handler) Create(ctx context.Context, r *pbmes.CreateReq) (*pbmes.CreateRes, error) {
	err := handler.Service.Create(r.UserId, r.RoomId, r.Text)
	if err != nil {
		return nil, err
	}
	return &pbmes.CreateRes{}, nil
}
func (handler *Handler) Delete(ctx context.Context, r *pbmes.DeleteReq) (*pbmes.DeleteRes, error) {
	err := handler.Service.Delete(r.MessageId, r.UserId)
	if err != nil {
		return nil, err
	}
	return &pbmes.DeleteRes{}, nil
}
func (handler *Handler) UpdateStatus(ctx context.Context, r *pbmes.UpdateStatusReq) (*pbmes.UpdateStatusRes, error) {
	err := handler.Service.UpdateStatus(r.MessageId, r.RoomId, r.UserId)
	if err != nil {
		return nil, err
	}
	return &pbmes.UpdateStatusRes{}, nil
}
