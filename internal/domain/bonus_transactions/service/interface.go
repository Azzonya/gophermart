package service

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonus_transactions/model"
)

type RepoAccrualI interface {
}

type RepoDbI interface {
	List(ctx context.Context, pars *model.ListPars) ([]*model.BonusTransaction, error)
	Create(ctx context.Context, obj *model.GetPars) error
	Get(ctx context.Context, pars *model.GetPars) (*model.BonusTransaction, bool, error)
	Update(ctx context.Context, pars *model.GetPars, obj *model.GetPars) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, login string) (bool, error)
}
