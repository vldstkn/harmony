package handlers

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"harmony/internal/config"
	"harmony/internal/models"
	"harmony/internal/services/api/middleware"
	"harmony/internal/services/ws/dto"
	pb "harmony/pkg/api/room"
	grpc_conn "harmony/pkg/grpc-conn"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RoomClient struct {
	Id     int64
	Conn   *websocket.Conn
	RoomId int64
}

type Room struct {
	Clients map[int64]*RoomClient
	Mutex   *sync.Mutex
}

type ChatHandlerDeps struct {
	Config   *config.Config
	Logger   *slog.Logger
	Producer sarama.SyncProducer
}

type ChatHandler struct {
	Config     *config.Config
	Logger     *slog.Logger
	Upgrader   *websocket.Upgrader
	RoomsMutex *sync.Mutex
	Rooms      map[int]*Room
	RoomClient pb.RoomClient
	Producer   sarama.SyncProducer
}

func NewChatHandler(router *chi.Mux, deps *ChatHandlerDeps) {
	roomConn, err := grpc_conn.NewClientConn(deps.Config.RoomAddress)
	if err != nil {
		panic(err)
	}
	roomClient := pb.NewRoomClient(roomConn)
	handler := &ChatHandler{
		Config:     deps.Config,
		Logger:     deps.Logger,
		Rooms:      make(map[int]*Room),
		RoomClient: roomClient,
		Upgrader:   &websocket.Upgrader{},
		RoomsMutex: &sync.Mutex{},
		Producer:   deps.Producer,
	}
	router.Route("/ch at", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.JWTSecret))
		r.Get("/ws/{room_id}", handler.WS())
	})
}

func (handler *ChatHandler) WS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomId, err := strconv.Atoi(chi.URLParam(r, "room_id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			handler.Logger.Error("chi.URLParam",
				slog.String("err", err.Error()),
				slog.String("Api addr", r.RemoteAddr))
			return
		}
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		response, _ := handler.RoomClient.CanSendMessage(context.Background(), &pb.CanSendMessageReq{
			RoomId: int64(roomId),
			UserId: userId,
		})
		if !response.CanSendMessage {
			handler.Logger.Error("RoomClient.CanSendMessage",
				slog.String("Api addr", r.RemoteAddr),
				slog.Int64("User id", userId),
				slog.Int("Room id", roomId),
			)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		conn, err := handler.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			handler.Logger.Error("Upgrader.Upgrade", slog.String("err", err.Error()))
			return
		}
		defer conn.Close()
		handler.Logger.Info("Client connect",
			slog.Int64("UserId", userId),
			slog.Int("RoomId", roomId))
		room := handler.getOrCreateRoom(roomId)
		client := &RoomClient{Id: userId, Conn: conn, RoomId: int64(roomId)}
		room.Mutex.Lock()
		room.Clients[userId] = client
		room.Mutex.Unlock()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			msgProd := dto.MessageProd{
				RoomId:   int64(roomId),
				SenderId: userId,
				Message:  string(msg),
			}
			msgProdByte, err := json.Marshal(msgProd)
			if err != nil {
				handler.Logger.Error(err.Error(),
					slog.String("Message", string(msg)),
					slog.Int64("User id", userId),
					slog.Int("Room id", roomId),
					slog.String("WS name", "Chat"),
				)
				continue
			}
			_, _, err = handler.Producer.SendMessage(&sarama.ProducerMessage{
				Topic: "message_create",
				Value: sarama.ByteEncoder(msgProdByte),
			})
			if err != nil {
				handler.Logger.Error("error when sending a message",
					slog.String("Message", string(msg)),
					slog.Int64("User id", userId),
				)
				continue
			}
			status := models.Unread
			msgJson := dto.MessageRes{
				Message:  string(msg),
				Status:   string(status),
				SenderId: userId,
				Date:     time.Now().UTC().String(),
			}
			msg, err = json.Marshal(msgJson)
			if err != nil {
				continue
			}
			room.broadcast(userId, msg)
		}
		room.Mutex.Lock()
		delete(room.Clients, userId)
		room.Mutex.Unlock()

		handler.Logger.Info("Client disconnect",
			slog.Int64("UserId", userId),
			slog.Int("RoomId", roomId),
			slog.String("WS name", "Chat"))
	}
}

func (room *Room) broadcast(senderId int64, msg []byte) {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()
	for id, client := range room.Clients {
		if id == senderId {
			continue
		}
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Conn.Close()
			delete(room.Clients, id)
		}
	}
}

func (handler *ChatHandler) getOrCreateRoom(chatId int) *Room {
	handler.RoomsMutex.Lock()
	defer handler.RoomsMutex.Unlock()

	if room, exists := handler.Rooms[chatId]; exists {
		return room
	}
	newRoom := &Room{
		Clients: make(map[int64]*RoomClient),
		Mutex:   &sync.Mutex{},
	}
	handler.Rooms[chatId] = newRoom
	return newRoom
}
