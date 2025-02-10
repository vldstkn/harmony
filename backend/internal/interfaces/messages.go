package interfaces

import "harmony/internal/models"

type MessageService interface {
	Create(userId, roomId int64, text string) error
	Delete(messageId, userId int64) error
	UpdateStatus(messageId, roomId, userId int64) error
	GetMessageByRoomId(roomId, startId, size int64) []models.Message
}
type MessageRepository interface {
	Create(userId, roomId int64, text string) error
	Delete(messageId int64) error
	UpdateStatus(lastMessage, nowMessage string) error
	GetMessageById(messageId int64) *models.Message
	GetMessageByRoomId(roomId, startId, size int64) []models.Message
	GetFirstUnreadMessage(roomId, userId int64) *models.Message
}
