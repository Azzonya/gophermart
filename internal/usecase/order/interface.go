package order

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/order"
)

type OrderServiceI interface {
	IsLuhnValid(orderNumber string) bool
	List(ctx context.Context, pars *order.ListPars) ([]*order.Order, error)
	ListWithAccrual(ctx context.Context, pars *order.ListPars) ([]*order.OrderWithAccrual, error)
	Get(ctx context.Context, pars *order.GetPars) (*order.Order, bool, error)
	Create(ctx context.Context, obj *order.GetPars) error
	Update(ctx context.Context, pars *order.GetPars) error
	Delete(ctx context.Context, pars *order.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
