package service

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
)

type Service struct {
	repoDb RepoDbI
}

func (s *Service) IsLoginTaken(ctx context.Context, login string) (bool, error) {
	return s.repoDb.Exists(ctx, login)
}

func (s *Service) Register(user *model.User) error {
	return nil
}
