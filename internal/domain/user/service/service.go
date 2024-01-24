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

func (s *Service) Register(ctx context.Context, user *model.GetPars) error {
	return s.repoDb.Create(ctx, user)
}

func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.User, error) {
	return s.repoDb.List(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *model.GetPars) error {
	return s.repoDb.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	return s.repoDb.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *model.GetPars, obj *model.GetPars) error {
	return s.repoDb.Update(ctx, pars, obj)
}

func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDb.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDb.Exists(ctx, orderNumber)
}
