package dto

import pb "harmony/pkg/api/account"

type AccountLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AccountLoginRes struct {
	Id          int64  `json:"id"`
	AccessToken string `json:"access_token"`
}

type AccountRegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

type AccountConfirmEmailReq struct {
	Token string `json:"token" validate:"required"`
}
type AccountConfirmEmailRes struct {
	Id          int64  `json:"id"`
	AccessToken string `json:"access_token"`
}
type FindByNameReq struct {
	Name string `json:"name" validate:"required"`
}
type FindByNameRes struct {
	Users []*pb.UserPublic `json:"users"`
}
type AddFriendReq struct {
	FriendId int64 `json:"friend_id" validate:"required"`
}
type DeleteFriendReq struct {
	FriendId int64 `json:"friend_id" validate:"required"`
}
type FindFriendsByNameReq struct {
	Name string `json:"name" validate:"required"`
}
type FindFriendsByNameRes struct {
	Users []*pb.UserPublic `json:"users"`
}
