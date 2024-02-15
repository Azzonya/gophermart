package bonustransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
)

type Service struct {
	repoDB BonusTransactionsRepoDBI
}

func New(repoDB BonusTransactionsRepoDBI) *Service {
	return &Service{
		repoDB: repoDB,
	}
}

func (s *Service) List(ctx context.Context, pars *entities.BonusTransactionsListPars) ([]*entities.BonusTransaction, error) {
	return s.repoDB.ListBt(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *entities.BonusTransaction) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *entities.BonusTransactionsParameters) (*entities.BonusTransaction, error) {
	return s.repoDB.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *entities.BonusTransactionsParameters) error {
	return s.repoDB.Update(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, pars *entities.BonusTransactionsParameters) error {
	return s.repoDB.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDB.Exists(ctx, orderNumber)
}
