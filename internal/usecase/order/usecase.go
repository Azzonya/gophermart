package order

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/order"
)

type Usecase struct {
	srv OrderServiceI
}

func New(srv OrderServiceI) *Usecase {
	return &Usecase{
		srv: srv,
	}
}

func (u *Usecase) IsLuhnValid(orderNumber string) bool {
	return u.srv.IsLuhnValid(orderNumber)
}

func (u *Usecase) List(ctx context.Context, pars *order.ListPars) ([]*order.Order, error) {
	return u.srv.List(ctx, pars)
}

func (u *Usecase) ListWithAccrual(ctx context.Context, pars *order.ListPars) ([]*order.OrderWithAccrual, error) {
	return u.srv.ListWithAccrual(ctx, pars)
}

func (u *Usecase) Get(ctx context.Context, pars *order.GetPars) (*order.Order, bool, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) Create(ctx context.Context, obj *order.GetPars) error {
	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Update(ctx context.Context, pars *order.GetPars) error {
	return u.srv.Update(ctx, pars)
}

func (u *Usecase) Delete(ctx context.Context, pars *order.GetPars) error {
	return u.srv.Delete(ctx, pars)
}

func (u *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return u.srv.Exists(ctx, orderNumber)
}
