package models

type MessageStatus string

const (
	Read   MessageStatus = "READ"
	Unread MessageStatus = "UNREAD"
)

type Message struct {
	Id        int64         `db:"id"`
	RoomId    int64         `db:"room_id"`
	SenderId  int64         `db:"sender_id"`
	Text      string        `db:"text"`
	CreatedAt string        `db:"created_at"`
	Status    MessageStatus `db:"status"`
}
