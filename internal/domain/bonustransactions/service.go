package bonustransactions

import (
	"context"
)

type Service struct {
	repoDB RepoDBIntreface
}

func New(repoDB RepoDBIntreface) *Service {
	return &Service{
		repoDB: repoDB,
	}
}

func (s *Service) List(ctx context.Context, pars *ListPars) ([]*BonusTransaction, error) {
	return s.repoDB.List(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *BonusTransaction) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *GetPars) (*BonusTransaction, error) {
	return s.repoDB.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *GetPars) error {
	return s.repoDB.Update(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, pars *GetPars) error {
	return s.repoDB.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDB.Exists(ctx, orderNumber)
}
