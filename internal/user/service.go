package user

import (
	"context"
	"github.com/commedesvlados/go-app/pkg/logging"
)

type Service struct {
	storage Storage
	logger  logging.Logger
}

func (s *Service) Create(ctx context.Context, dto CreateUserDTO) (u User, err error) {
	return
}
