package interfaces

type ChatService interface {
	SaveMessage(userId, roomId int64, message []byte) (string, error)
}

type ChatRepository interface {
	SaveMessage(userId, roomId int64, message string) (string, error)
}
