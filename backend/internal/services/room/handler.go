package room

import (
	"context"
	"errors"
	"harmony/internal/config"
	"harmony/internal/interfaces"
	pb "harmony/pkg/api/room"
	"log/slog"
	"net/http"
)

type HandlerDeps struct {
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.RoomService
}
type Handler struct {
	pb.UnsafeRoomServer
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.RoomService
}

func NewHandler(deps *HandlerDeps) *Handler {
	return &Handler{
		Config:  deps.Config,
		Logger:  deps.Logger,
		Service: deps.Service,
	}
}

func (handler *Handler) CreateRoom(ctx context.Context, r *pb.CreateRoomReq) (*pb.CreateRoomRes, error) {
	roomId, err := handler.Service.CreateRoom(r.CreatorId, r.Name, r.Participants)
	if err != nil {
		return nil, err
	}
	return &pb.CreateRoomRes{
		RoomId: roomId,
	}, nil
}
func (handler *Handler) DeleteRoom(ctx context.Context, r *pb.DeleteRoomReq) (*pb.DeleteRoomRes, error) {
	err := handler.Service.DeleteRoom(r.UserId, r.RoomId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteRoomRes{}, nil
}
func (handler *Handler) AddUsersToRoom(ctx context.Context, r *pb.AddUsersToRoomReq) (*pb.AddUsersToRoomRes, error) {
	err := handler.Service.AddUsersToRoom(r.CreatorId, r.RoomId, r.UsersId)
	if err != nil {
		return nil, err
	}
	return &pb.AddUsersToRoomRes{}, nil
}
func (handler *Handler) RemoveUsersFromRoom(ctx context.Context, r *pb.RemoveUsersFromRoomReq) (*pb.RemoveUsersFromRoomRes, error) {
	err := handler.Service.RemoveUsersFromRoom(r.CreatorId, r.RoomId, r.UsersId)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveUsersFromRoomRes{}, nil
}
func (handler *Handler) GetRoomsByUserId(ctx context.Context, r *pb.GetRoomsByUserIdReq) (*pb.GetRoomsByUserIdRes, error) {
	rooms := handler.Service.GetRoomsByUserId(r.UserId)
	return &pb.GetRoomsByUserIdRes{
		Rooms: FromModelRoomsToPublic(rooms),
	}, nil
}

func (handler *Handler) GetRoom(ctx context.Context, r *pb.GetRoomReq) (*pb.GetRoomRes, error) {
	room := handler.Service.GetRoomById(r.UserId, r.RoomId)
	if room == nil {
		return nil, errors.New(http.StatusText(http.StatusForbidden))
	}
	return &pb.GetRoomRes{
		Room: FromModelRoomToPublic(*room),
	}, nil
}
func (handler *Handler) CanSendMessage(ctx context.Context, r *pb.CanSendMessageReq) (*pb.CanSendMessageRes, error) {
	canSend := handler.Service.CanSendMessage(r.UserId, r.RoomId)
	return &pb.CanSendMessageRes{
		CanSendMessage: canSend,
	}, nil
}
