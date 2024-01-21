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

func (s *Service) Register(ctx context.Context, user *model.User) error {
	return s.repoDb.Create(ctx, user)
}

func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	return s.repoDb.Get(ctx, pars)
}
