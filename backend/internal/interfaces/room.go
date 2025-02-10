package interfaces

import "harmony/internal/models"

type RoomService interface {
	CreateRoom(creatorId int64, name string, usersId []int64) (int64, error)
	DeleteRoom(creatorId, roomId int64) error
	AddUsersToRoom(creatorId, roomId int64, usersId []int64) error
	RemoveUsersFromRoom(creatorId, roomId int64, usersId []int64) error
	GetRoomsByUserId(userId int64) []models.Room
	GetRoomById(userId, roomId int64) *models.Room
	CanSendMessage(userId, roomId int64) bool
}

type RoomRepository interface {
	Create(creatorId int64, name string) (int64, error)
	Delete(roomId int64) error
	AddUsers(roomId int64, usersId []int64) error
	RemoveUsers(roomId int64, usersId []int64) error
	GetRoomsByUserId(userId int64) []models.Room
	GetRoomById(id int64) *models.Room
	GetRoomParticipants(roomId int64) []int64
	CheckAndGetRoomForUser(roomId, userId int64) *models.Room
	GetRoomRole(roomId, userId int64) models.RoomRole
	GetRoomRoles(roomId int64, usersId []int64) []models.RoomMember
	GetRoomMember(userId, roomId int64) *models.RoomMember
}
