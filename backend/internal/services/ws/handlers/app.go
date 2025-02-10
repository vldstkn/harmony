package handlers

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"harmony/internal/config"
	"harmony/internal/contracts"
	"harmony/internal/services/api/middleware"
	"log/slog"
	"net/http"
	"sync"
)

type AppClient struct {
	Id   int64
	Conn *websocket.Conn
}

type AppHandlerDeps struct {
	Config    *config.Config
	Logger    *slog.Logger
	Producer  sarama.SyncProducer
	ChanKafka chan []byte
}

type AppHandler struct {
	Config    *config.Config
	Logger    *slog.Logger
	Upgrader  *websocket.Upgrader
	Producer  sarama.SyncProducer
	Clients   map[int64]*AppClient
	Mutex     *sync.Mutex
	ChanKafka chan []byte
}

func NewAppHandler(router *chi.Mux, deps *AppHandlerDeps) *AppHandler {
	handler := &AppHandler{
		Config:    deps.Config,
		Logger:    deps.Logger,
		Upgrader:  &websocket.Upgrader{},
		Producer:  deps.Producer,
		Mutex:     &sync.Mutex{},
		Clients:   make(map[int64]*AppClient),
		ChanKafka: deps.ChanKafka,
	}
	router.Route("/app", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.JWTSecret))
		r.Get("/ws", handler.WS())
	})

	return handler
}

func (handler *AppHandler) WS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		conn, err := handler.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			handler.Logger.Error("Upgrader.Upgrade", slog.String("err", err.Error()))
			return
		}
		defer conn.Close()
		handler.Logger.Info("Client connect",
			slog.Int64("UserId", userId),
			slog.String("IP", r.RemoteAddr),
			slog.String("WS name", "App"))
		handler.Mutex.Lock()
		handler.Clients[userId] = &AppClient{
			Id:   userId,
			Conn: conn,
		}
		handler.Mutex.Unlock()
		go handler.BroadcastMessage()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
		handler.Logger.Info("Client disconnect",
			slog.Int64("UserId", userId),
			slog.String("IP", r.RemoteAddr),
			slog.String("WS name", "App"),
		)
		handler.Mutex.Lock()
		delete(handler.Clients, userId)
		handler.Mutex.Unlock()
	}
}

func (handler *AppHandler) BroadcastMessage() {
	for message := range handler.ChanKafka {
		var mes contracts.Notification
		err := json.Unmarshal(message, &mes)
		if err != nil {
			handler.Logger.Error(err.Error(),
				slog.String("Error Loacation", "json.Unmarshal"))
			continue
		}
		handler.Mutex.Lock()
		if client, exists := handler.Clients[mes.UserId]; exists {
			err := client.Conn.WriteJSON(mes)
			if err != nil {
				client.Conn.Close()
				delete(handler.Clients, mes.UserId)
			}
		}
		handler.Mutex.Unlock()
	}
}
