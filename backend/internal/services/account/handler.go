package account

import (
	"context"
	"errors"
	"harmony/internal/config"
	"harmony/internal/interfaces"
	pb "harmony/pkg/api/account"
	pbE "harmony/pkg/api/email"
	grpc_conn "harmony/pkg/grpc-conn"
	"harmony/pkg/jwt"
	"log/slog"
)

type Handler struct {
	pb.UnsafeAccountServer
	Config      *config.Config
	Logger      *slog.Logger
	Service     interfaces.AccountService
	EmailClient pbE.EmailClient
}

type HandlerDeps struct {
	Config  *config.Config
	Logger  *slog.Logger
	Service interfaces.AccountService
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
