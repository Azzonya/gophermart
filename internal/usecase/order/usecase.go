package order

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/order/model"
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

func (s *Usecase) List(ctx context.Context, pars *model.ListPars) ([]*model.Order, error) {
	return s.srv.List(ctx, pars)
}

func (u *Usecase) Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) Create(ctx context.Context, obj *model.GetPars) error {
	return u.srv.Create(ctx, obj)
}

func (s *Usecase) Update(ctx context.Context, pars *model.GetPars, obj *model.GetPars) error {
	return s.srv.Update(ctx, pars, obj)
}

func (s *Usecase) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.srv.Delete(ctx, pars)
}

func (s *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.srv.Exists(ctx, orderNumber)
}
