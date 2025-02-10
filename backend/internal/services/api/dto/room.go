package dto

import pb "harmony/pkg/api/room"

type RoomCreateReq struct {
	Name           string  `json:"name" validate:"required"`
	ParticipantsId []int64 `json:"participants_id,omitempty" validate:"omitnil,min=1,dive,gte=0"`
}

type RoomCreateRes struct {
	RoomId int64 `json:"room_id"`
}

type RoomDeleteReq struct {
	RoomId int64 `json:"room_id" validate:"required,number"`
}

type AddUsersToRoomReq struct {
	RoomId         int64   `json:"room_id" validate:"required,number"`
	ParticipantsId []int64 `json:"participants_id,omitempty" validate:"required,min=1,dive,gte=0"`
}
type RemoveUsersFromRoomReq struct {
	RoomId         int64   `json:"room_id" validate:"required,number"`
	ParticipantsId []int64 `json:"participants_id,omitempty" validate:"required,min=1,dive,gte=0"`
}
type GetRoomsByUserIdRes struct {
	Rooms []*pb.RoomPublic
}
type GetRoomRes struct {
	Room pb.RoomPublic `json:"room"`
}
