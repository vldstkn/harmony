package dto

type MessageRes struct {
	SenderId int64
	Message  string
	Date     string
	Status   string
}

type MessageProd struct {
	SenderId int64  `json:"user_id"`
	RoomId   int64  `json:"room_id"`
	Message  string `json:"message"`
}

type MessageNotificationReq struct {
	UsersId  []int64
	RoomId   int64
	Text     string
	Date     string
	SenderId int64
}
