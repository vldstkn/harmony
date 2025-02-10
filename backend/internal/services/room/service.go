package room

import (
	"errors"
	"harmony/internal/interfaces"
	"harmony/internal/models"
	"log/slog"
	"net/http"
)

type ServiceDeps struct {
	Logger     *slog.Logger
	Repository interfaces.RoomRepository
}
type Service struct {
	Logger     *slog.Logger
	Repository interfaces.RoomRepository
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Logger:     deps.Logger,
		Repository: deps.Repository,
	}
}

func (service *Service) CreateRoom(creatorId int64, name string, usersId []int64) (int64, error) {
	roomId, err := service.Repository.Create(creatorId, name)
	if err != nil {
		service.Logger.Error("Repository.Create", slog.String("err", err.Error()),
			slog.Int64("creatorId", creatorId),
			slog.String("name", name))
		return -1, errors.New(http.StatusText(http.StatusInternalServerError))
	}
	if len(usersId) > 0 {
		err = service.Repository.AddUsers(roomId, usersId)
		if err != nil {
			service.Logger.Error("Repository.AddUsers", slog.String("err", err.Error()),
				slog.Int64("roomId", roomId))
			return -1, errors.New(http.StatusText(http.StatusInternalServerError))
		}
	}
	return roomId, nil
}
func (service *Service) DeleteRoom(creatorId, roomId int64) error {
	role := service.Repository.GetRoomRole(roomId, creatorId)
	if role != models.Owner {
		return errors.New(http.StatusText(http.StatusForbidden))
	}
	err := service.Repository.Delete(roomId)
	if err != nil {
		service.Logger.Error("Repository.Delete", slog.String("err", err.Error()),
			slog.Int64("roomId", roomId),
			slog.Int64("creatorId", creatorId))
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
func (service *Service) AddUsersToRoom(creatorId, roomId int64, usersId []int64) error {
	role := service.Repository.GetRoomRole(roomId, creatorId)
	if role < 1 {
		return errors.New(http.StatusText(http.StatusForbidden))
	}
	err := service.Repository.AddUsers(roomId, usersId)
	if err != nil {
		service.Logger.Error("Repository.AddUsers", slog.String("err", err.Error()),
			slog.Int64("roomId", roomId))
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
func (service *Service) RemoveUsersFromRoom(creatorId, roomId int64, usersId []int64) error {
	role := service.Repository.GetRoomRole(roomId, creatorId)
	if role < 1 {
		return errors.New(http.StatusText(http.StatusForbidden))
	}
	err := service.Repository.RemoveUsers(roomId, usersId)
	if err != nil {
		service.Logger.Error("Repository.RemoveUsers", slog.String("err", err.Error()),
			slog.Int64("roomId", roomId),
			slog.Int64("creatorId", creatorId),
		)
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}
func (service *Service) GetRoomsByUserId(userId int64) []models.Room {
	rooms := service.Repository.GetRoomsByUserId(userId)
	return rooms
}

func (service *Service) GetRoomById(userId, roomId int64) *models.Room {
	room := service.Repository.CheckAndGetRoomForUser(roomId, userId)
	return room
}
func (service *Service) CanSendMessage(userId, roomId int64) bool {
	roomMember := service.Repository.GetRoomMember(userId, roomId)
	if roomMember == nil || roomMember.Role < 0 {
		return false
	}
	return true
}
