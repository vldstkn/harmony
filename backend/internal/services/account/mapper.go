package account

import (
	"harmony/internal/models"
	pb "harmony/pkg/api/account"
)

func UserFromModelToProto(user models.User) *pb.UserPublic {
	return &pb.UserPublic{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func UsersFromModelToProto(users []models.User) []*pb.UserPublic {
	pbUsers := make([]*pb.UserPublic, len(users))
	for i, el := range users {
		pbUsers[i] = UserFromModelToProto(el)
	}
	return pbUsers
}
