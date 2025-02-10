package notifications

import (
	"harmony/internal/interfaces"
	"log/slog"
)

type ServiceDeps struct {
	Logger     *slog.Logger
	Repository interfaces.NotificationsRepository
}

type Service struct {
	Logger     *slog.Logger
	Repository interfaces.NotificationsRepository
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Logger:     deps.Logger,
		Repository: deps.Repository,
	}
}
func (service *Service) Save() {
	return
}
