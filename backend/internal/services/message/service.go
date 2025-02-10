package messages

import (
	"errors"
	"harmony/internal/interfaces"
	"harmony/internal/models"
	"log/slog"
	"net/http"
)

type ServiceDeps struct {
	Repository interfaces.MessageRepository
	Logger     *slog.Logger
}
type Service struct {
	Repository interfaces.MessageRepository
	Logger     *slog.Logger
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Repository: deps.Repository,
		Logger:     deps.Logger,
	}
}
func (service *Service) Create(userId, roomId int64, text string) error {
	err := service.Repository.Create(userId, roomId, text)
	if err != nil {
		return errors.New(http.StatusText(http.StatusBadRequest))
	}
	return nil
}
func (service *Service) Delete(messageId, userId int64) error {
	message := service.Repository.GetMessageById(messageId)
	if message == nil {
		service.Logger.Error("message not found",
			slog.String("Error location", "Repository.GetMessageById"),
			slog.Int64("Message id", messageId))
		return errors.New(http.StatusText(http.StatusBadRequest))
	}
	if message.SenderId != userId {
		service.Logger.Error("the sender's ID does not match the user's ID",
			slog.String("Error location", "Repository.GetMessageById"),
			slog.Int64("User id", userId),
			slog.Int64("Sender id", message.SenderId),
			slog.Int64("Message id", messageId))
		return errors.New(http.StatusText(http.StatusForbidden))
	}
	err := service.Repository.Delete(messageId)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", "Repository.Delete"),
			slog.Int64("Message id", messageId))
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
func (service *Service) UpdateStatus(messageId, roomId, userId int64) error {
	nowMessage := service.Repository.GetMessageById(messageId)
	if nowMessage == nil {
		service.Logger.Error("message not found",
			slog.String("Error location", "Repository.UpdateStatus"),
			slog.Int64("Message id", messageId),
		)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}
	lastMessage := service.Repository.GetFirstUnreadMessage(roomId, userId)
	if lastMessage == nil {
		lastMessage = nowMessage
	}
	err := service.Repository.UpdateStatus(lastMessage.CreatedAt, nowMessage.CreatedAt)
	if err != nil {
		service.Logger.Error(err.Error(),
			slog.String("Error location", "Repository.UpdateStatus"),
			slog.Int64("Last message id", lastMessage.Id),
			slog.String("Last message date", lastMessage.CreatedAt),
			slog.Int64("Now message id", nowMessage.Id),
			slog.String("Now message date", nowMessage.CreatedAt),
		)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}
	return nil
}
func (service *Service) GetMessageByRoomId(roomId, startId, size int64) []models.Message {
	messages := service.Repository.GetMessageByRoomId(roomId, startId, size)
	return messages
}
