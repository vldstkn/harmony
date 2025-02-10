package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"harmony/internal/config"
	"harmony/internal/services/api/dto"
	"harmony/internal/services/api/middleware"
	pb "harmony/pkg/api/room"
	grpc_conn "harmony/pkg/grpc-conn"
	"harmony/pkg/req"
	"harmony/pkg/res"
	"log/slog"
	"net/http"
	"strconv"
)

type RoomHandlerDeps struct {
	Config *config.Config
	Logger *slog.Logger
}

type RoomHandler struct {
	Config     *config.Config
	RoomClient pb.RoomClient
	Logger     *slog.Logger
}

func NewRoomHandler(router *chi.Mux, deps *RoomHandlerDeps) {
	roomConn, err := grpc_conn.NewClientConn(deps.Config.RoomAddress)
	if err != nil {
		panic(err)
	}
	roomClient := pb.NewRoomClient(roomConn)
	handler := &RoomHandler{
		Config:     deps.Config,
		Logger:     deps.Logger,
		RoomClient: roomClient,
	}

	router.Route("/rooms", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.JWTSecret))
		r.Post("/", handler.CreateRoom())
		r.Post("/add-users", handler.AddUsersToRoom())
		r.Post("/remove-users", handler.RemoveUsersFromRoom())

		r.Get("/", handler.GetRoomsByUserId())
		r.Get("/{id}", handler.GetRoom())

		r.Delete("/", handler.DeleteRoom())
	})
}

func (handler *RoomHandler) CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.RoomCreateReq](r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		creatorId := r.Context().Value("authData").(middleware.AuthData).Id
		response, err := handler.RoomClient.CreateRoom(context.Background(), &pb.CreateRoomReq{
			Participants: body.ParticipantsId,
			Name:         body.Name,
			CreatorId:    creatorId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, dto.RoomCreateRes{
			RoomId: response.RoomId,
		}, http.StatusCreated)
	}
}
func (handler *RoomHandler) DeleteRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.RoomDeleteReq](r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		_, err = handler.RoomClient.DeleteRoom(context.Background(), &pb.DeleteRoomReq{
			RoomId: body.RoomId,
			UserId: userId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, nil, http.StatusOK)
	}
}
func (handler *RoomHandler) AddUsersToRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.AddUsersToRoomReq](r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		_, err = handler.RoomClient.AddUsersToRoom(context.Background(), &pb.AddUsersToRoomReq{
			RoomId:    body.RoomId,
			UsersId:   body.ParticipantsId,
			CreatorId: userId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, nil, http.StatusOK)
	}
}
func (handler *RoomHandler) RemoveUsersFromRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.RemoveUsersFromRoomReq](r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		_, err = handler.RoomClient.RemoveUsersFromRoom(context.Background(), &pb.RemoveUsersFromRoomReq{
			RoomId:    body.RoomId,
			UsersId:   body.ParticipantsId,
			CreatorId: userId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, nil, http.StatusOK)
	}
}
func (handler *RoomHandler) GetRoomsByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		response, _ := handler.RoomClient.GetRoomsByUserId(context.Background(), &pb.GetRoomsByUserIdReq{
			UserId: userId,
		})
		res.Json(w, dto.GetRoomsByUserIdRes{
			Rooms: response.Rooms,
		}, http.StatusOK)
	}
}
func (handler *RoomHandler) GetRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			handler.Logger.Error("chi.URLParam", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		response, err := handler.RoomClient.GetRoom(context.Background(), &pb.GetRoomReq{
			UserId: userId,
			RoomId: int64(roomId),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, dto.GetRoomRes{
			Room: *response.Room,
		}, http.StatusOK)
	}
}
