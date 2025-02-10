package messages

type MessageReq struct {
	SenderId int64  `json:"user_id"`
	RoomId   int64  `json:"room_id"`
	Message  string `json:"message"`
}
