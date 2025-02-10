package account

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"harmony/internal/config"
	"harmony/internal/contracts"
	"harmony/internal/interfaces"
	pb "harmony/pkg/api/account"
	pbE "harmony/pkg/api/email"
	grpc_conn "harmony/pkg/grpc-conn"
	"harmony/pkg/jwt"
	"log/slog"
	"time"
)

type Handler struct {
	pb.UnsafeAccountServer
	Config      *config.Config
	Logger      *slog.Logger
	Service     interfaces.AccountService
	EmailClient pbE.EmailClient
	Producer    sarama.SyncProducer
}

type HandlerDeps struct {
	Config   *config.Config
	Logger   *slog.Logger
	Service  interfaces.AccountService
	Producer sarama.SyncProducer
}

func NewHandler(deps *HandlerDeps) *Handler {
	emailConn, err := grpc_conn.NewClientConn(deps.Config.EmailAddress)
	if err != nil {
		panic(err)
	}
	emailClient := pbE.NewEmailClient(emailConn)
	return &Handler{
		Config:      deps.Config,
		Logger:      deps.Logger,
		Service:     deps.Service,
		EmailClient: emailClient,
		Producer:    deps.Producer,
	}
}

func (handler *Handler) Register(ctx context.Context, r *pb.RegisterReq) (*pb.RegisterRes, error) {
	id, err := handler.Service.Register(r.Email, r.Password, r.Name)
	if err != nil {
		handler.Logger.Error(err.Error(), slog.String("email", r.Email), slog.String("name", r.Name))
		return nil, err
	}
	if err != nil {
		handler.Logger.Error(err.Error())
		return nil, err
	}
	token, _, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id: id,
	})
	go handler.EmailClient.ConfirmEmail(context.Background(), &pbE.ConfirmEmailSendReq{
		Email: r.Email,
		Token: token,
	})

	return &pb.RegisterRes{
		Email: r.Email,
	}, nil

}
func (handler *Handler) Login(ctx context.Context, r *pb.LoginReq) (*pb.LoginRes, error) {
	id, err := handler.Service.Login(r.Email, r.Password)
	if err != nil {
		return nil, err
	}
	accT, refT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id: id,
	})
	if err != nil {
		return nil, errors.New("internal error")
	}
	return &pb.LoginRes{
		Id:           id,
		RefreshToken: accT,
		AccessToken:  refT,
	}, nil
}
func (handler *Handler) ConfirmEmail(ctx context.Context, r *pb.ConfirmEmailReq) (*pb.ConfirmEmailRes, error) {
	id, err := handler.Service.ConfirmEmail(handler.Config.JWTSecret, r.Token)
	if err != nil {
		return nil, err
	}
	accT, refT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id: id,
	})
	if err != nil {
		return nil, errors.New("internal error")
	}
	return &pb.ConfirmEmailRes{
		AccessToken:  accT,
		RefreshToken: refT,
		Id:           id,
	}, nil
}

func (handler *Handler) GetNewTokens(ctx context.Context, r *pb.GetNewTokensReq) (*pb.GetNewTokensRes, error) {
	isValid, data := jwt.NewJWT(handler.Config.JWTSecret).Parse(r.RefreshToken)
	if !isValid {
		handler.Logger.Error("token not valid")
		return nil, errors.New("token not valid")
	}
	aT, rT, err := handler.Service.IssueTokens(handler.Config.JWTSecret, jwt.Data{
		Id: data.Id,
	})
	if err != nil {
		return nil, errors.New("internal error")
	}
	return &pb.GetNewTokensRes{
		RefreshToken: rT,
		AccessToken:  aT,
	}, nil
}
func (handler *Handler) FindByName(ctx context.Context, r *pb.FindByNameReq) (*pb.FindByNameRes, error) {
	users := handler.Service.FindByName(r.Id, r.Name)
	res := UsersFromModelToProto(users)
	return &pb.FindByNameRes{
		Users: res,
	}, nil
}
func (handler *Handler) AddFriend(ctx context.Context, r *pb.AddFriendReq) (*pb.AddFriendRes, error) {
	//err := handler.Service.AddFriend(r.UserId, r.FriendId)
	// FIXME: вернуть все как было
	if true {
		message := contracts.Notification{
			Date:        time.Now().UTC().String(),
			Description: fmt.Sprintf("Вас хочет добавить в друзья пользователь %s", r.UserId),
			UserId:      r.FriendId,
			Title:       "Запрос в друзья",
			Type:        "FRIEND",
		}
		msgByte, _ := json.Marshal(message)
		_, _, err := handler.Producer.SendMessage(&sarama.ProducerMessage{
			Value: sarama.ByteEncoder(msgByte),
			Topic: "add_friend",
		})
		if err != nil {
			handler.Logger.Error(err.Error(),
				slog.String("Error location", "Producer.SendMessage"),
				slog.Int64("User id", r.UserId),
				slog.Int64("Friend id", r.FriendId),
			)
		}
	}
	return nil, nil
}
func (handler *Handler) DeleteFriend(ctx context.Context, r *pb.DeleteFriendReq) (*pb.DeleteFriendRes, error) {
	err := handler.Service.DeleteFriend(r.UserId, r.FriendId)
	return nil, err
}
func (handler *Handler) FindFriendsByName(ctx context.Context, r *pb.FindFriendsByNameReq) (*pb.FindFriendsByNameRes, error) {

	users := handler.Service.FindFriendsByName(r.UserId, r.Name)
	res := UsersFromModelToProto(users)
	return &pb.FindFriendsByNameRes{
		Users: res,
	}, nil
}
