package service

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/order/model"
)

type RepoDbI interface {
	List(ctx context.Context, pars *model.ListPars) ([]*model.Order, error)
	Create(ctx context.Context, obj *model.GetPars) error
	Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error)
	Update(ctx context.Context, pars *model.GetPars) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
