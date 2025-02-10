package messages

import (
	"harmony/internal/models"
	"harmony/pkg/db"
)

type Repository struct {
	DB *db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) Create(userId, roomId int64, text string) error {
	_, err := repo.DB.Exec(`INSERT INTO messages (room_id, sender_id, text) VALUES ($1,$2,$3)`,
		roomId, userId, text)
	if err != nil {
		return err
	}
	return nil
}
func (repo *Repository) Delete(messageId int64) error {
	_, err := repo.DB.Exec(`DELETE FROM messages WHERE id=$1`, messageId)
	return err
}
func (repo *Repository) UpdateStatus(lastMessage, nowMessage string) error {
	_, err := repo.DB.Exec(`UPDATE messages SET status='READ' WHERE created_at BETWEEN $1 AND $2`, lastMessage, nowMessage)
	return err
}
func (repo *Repository) GetMessageByRoomId(roomId, startId, size int64) []models.Message {
	var messages []models.Message
	err := repo.DB.Select(&messages, `SELECT * FROM messages 
																			     WHERE room_id=$1 AND id >= $2 
																			     LIMIT $3`, roomId, startId, size)
	if err != nil {
		return nil
	}
	return messages
}
func (repo *Repository) GetMessageById(messageId int64) *models.Message {
	var message models.Message
	err := repo.DB.Get(&message, `SELECT * FROM messages WHERE id=$1`, messageId)
	if err != nil {
		return nil
	}
	return &message
}

func (repo *Repository) GetFirstUnreadMessage(roomId, userId int64) *models.Message {
	var message models.Message
	err := repo.DB.Get(&message, `SELECT * FROM messages 
													      WHERE room_id=$1 AND status='UNREAD' AND sender_id!=$2
													      ORDER BY created_at LIMIT 1`, roomId, userId)
	if err != nil {
		return nil
	}
	return &message
}
