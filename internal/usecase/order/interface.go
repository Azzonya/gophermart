package order

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/order/model"
)

type OrderServiceI interface {
	IsLuhnValid(orderNumber string) bool
	List(ctx context.Context, pars *model.ListPars) ([]*model.Order, error)
	ListWithAccrual(ctx context.Context, pars *model.ListPars) ([]*model.OrderWithAccrual, error)
	Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error)
	Create(ctx context.Context, obj *model.GetPars) error
	Update(ctx context.Context, pars *model.GetPars) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
