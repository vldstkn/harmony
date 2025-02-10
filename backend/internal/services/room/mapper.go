package room

import (
	"harmony/internal/models"
	pb "harmony/pkg/api/room"
)

func FromModelRoomsToPublic(rooms []models.Room) []*pb.RoomPublic {
	pbRooms := make([]*pb.RoomPublic, len(rooms))
	for i := 0; i < len(rooms); i++ {
		pbRooms[i] = FromModelRoomToPublic(rooms[i])
	}
	return pbRooms
}
func FromModelRoomToPublic(room models.Room) *pb.RoomPublic {
	return &pb.RoomPublic{
		Id:             room.Id,
		CreatorId:      room.CreatorId,
		Name:           room.Name,
		ParticipantsId: room.Participants,
	}
}
