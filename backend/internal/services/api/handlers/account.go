package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"harmony/internal/config"
	"harmony/internal/interfaces"
	"harmony/internal/services/api/dto"
	"harmony/internal/services/api/middleware"
	pb "harmony/pkg/api/account"
	grpc_conn "harmony/pkg/grpc-conn"
	"harmony/pkg/req"
	"harmony/pkg/res"
	"log/slog"
	"net/http"
)

type AccountHandlerDeps struct {
	Service interfaces.ApiService
	Config  *config.Config
	Logger  *slog.Logger
}

type AccountHandler struct {
	Service       interfaces.ApiService
	Config        *config.Config
	AccountClient pb.AccountClient
	Logger        *slog.Logger
}

func NewAccountHandler(router *chi.Mux, deps *AccountHandlerDeps) {
	accountConn, err := grpc_conn.NewClientConn(deps.Config.AccountAddress)
	if err != nil {
		panic(err)
	}
	accountClient := pb.NewAccountClient(accountConn)
	handler := &AccountHandler{
		Service:       deps.Service,
		Config:        deps.Config,
		AccountClient: accountClient,
		Logger:        deps.Logger,
	}

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register())
		r.Post("/login", handler.Login())
		r.Post("/confirm-email", handler.ConfirmEmail())
	})

	router.Route("/users", func(r chi.Router) {
		r.Use(middleware.IsAuthed(handler.Config.JWTSecret))
		r.Get("/", handler.FindByName())
		r.Post("/add-friend", handler.AddFriend())
		r.Post("/delete-friend", handler.DeleteFriend())
		r.Get("/friends", handler.FindFriendsByName())
	})
}

func (handler *AccountHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.AccountRegisterReq](r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		response, err := handler.AccountClient.Register(context.Background(), &pb.RegisterReq{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		res.Json(w, response.Email, 200)
	}
}
func (handler *AccountHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.AccountLoginReq](r)
		if err != nil {
			handler.Logger.Error("req.HandleBody", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		response, err := handler.AccountClient.Login(context.Background(), &pb.LoginReq{
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			handler.Logger.Error("AccountClient.Login", slog.String("err", err.Error()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		handler.Service.AddCookie(&w, "refresh_token", response.RefreshToken, 3600)
		res.Json(w, dto.AccountLoginRes{
			Id:          response.Id,
			AccessToken: response.AccessToken,
		}, 200)
	}
}
func (handler *AccountHandler) ConfirmEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[dto.AccountConfirmEmailReq](r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		response, err := handler.AccountClient.ConfirmEmail(context.Background(), &pb.ConfirmEmailReq{
			Token: body.Token,
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		handler.Service.AddCookie(&w, "refresh_token", response.RefreshToken, 3600)
		res.Json(w, dto.AccountConfirmEmailRes{
			AccessToken: response.AccessToken,
			Id:          response.Id,
		}, 200)
	}
}

func (handler *AccountHandler) FindByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		name := r.URL.Query().Get("name")
		response, _ := handler.AccountClient.FindByName(context.Background(), &pb.FindByNameReq{
			Id:   userId,
			Name: name,
		})
		res.Json(w, dto.FindByNameRes{
			Users: response.Users,
		}, 200)
	}
}

func (handler *AccountHandler) AddFriend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		body, err := req.HandleBody[dto.AddFriendReq](r)
		if err != nil {
			handler.Logger.Error("req.HandleBody[dto.AddFriendReq]", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		_, err = handler.AccountClient.AddFriend(context.Background(), &pb.AddFriendReq{
			FriendId: body.FriendId,
			UserId:   userId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, nil, 201)
	}
}
func (handler *AccountHandler) DeleteFriend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		body, err := req.HandleBody[dto.DeleteFriendReq](r)
		if err != nil {
			handler.Logger.Error("req.HandleBody[dto.DeleteFriendReq]", slog.String("err", err.Error()))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		_, err = handler.AccountClient.DeleteFriend(context.Background(), &pb.DeleteFriendReq{
			FriendId: body.FriendId,
			UserId:   userId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, nil, 201)
	}
}
func (handler *AccountHandler) FindFriendsByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("authData").(middleware.AuthData).Id
		name := r.URL.Query().Get("name")
		response, _ := handler.AccountClient.FindFriendsByName(context.Background(), &pb.FindFriendsByNameReq{
			Name:   name,
			UserId: userId,
		})
		res.Json(w, dto.FindFriendsByNameRes{
			Users: response.Users,
		}, 200)
	}
}
