package service

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonus_transactions/model"
)

type Service struct {
	repoDb RepoDbI
}

func New(repoDb RepoDbI) *Service {
	return &Service{
		repoDb: repoDb,
	}
}

func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.BonusTransaction, error) {
	return s.repoDb.List(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *model.GetPars) error {
	return s.repoDb.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.BonusTransaction, bool, error) {
	return s.repoDb.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *model.GetPars) error {
	return s.repoDb.Update(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDb.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDb.Exists(ctx, orderNumber)
}
