package order

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
)

type OrderServiceI interface {
	IsLuhnValid(orderNumber string) bool
	List(ctx context.Context, pars *entities.OrderListPars) ([]*entities.Order, error)
	ListWithAccrual(ctx context.Context, pars *entities.OrderListPars) ([]*entities.OrderWithAccrual, error)
	Get(ctx context.Context, pars *entities.OrderParameters) (*entities.Order, error)
	Create(ctx context.Context, obj *entities.Order) error
	Update(ctx context.Context, pars *entities.OrderParameters) error
	Delete(ctx context.Context, pars *entities.OrderParameters) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
