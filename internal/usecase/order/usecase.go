package order

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/Azzonya/gophermart/internal/errs"
)

type Usecase struct {
	srv OrderServiceI
}

func New(srv OrderServiceI) *Usecase {
	return &Usecase{
		srv: srv,
	}
}

func (u *Usecase) List(ctx context.Context, pars *entities.OrderListPars) ([]*entities.Order, error) {
	return u.srv.List(ctx, pars)
}

func (u *Usecase) ListWithAccrual(ctx context.Context, pars *entities.OrderListPars) ([]*entities.OrderWithAccrual, error) {
	return u.srv.ListWithAccrual(ctx, pars)
}

func (u *Usecase) Get(ctx context.Context, pars *entities.OrderParameters) (*entities.Order, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) Create(ctx context.Context, obj *entities.Order) error {
	if !u.srv.IsLuhnValid(obj.OrderNumber) {
		return errs.ErrOrderNumberLuhnValid{OrderNumber: obj.OrderNumber}
	}

	foundOrder, err := u.Get(ctx, &entities.OrderParameters{
		OrderNumber: obj.OrderNumber,
	})
	if err != nil {
		return err
	}

	if foundOrder != nil {
		if obj.UserID == foundOrder.UserID {
			return errs.ErrOrderUploaded{OrderNumber: obj.OrderNumber}
		} else {
			return errs.ErrOrderUploadedByAnotherUser{OrderNumber: obj.OrderNumber}
		}
	}

	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Update(ctx context.Context, pars *entities.OrderParameters) error {
	return u.srv.Update(ctx, pars)
}

func (u *Usecase) Delete(ctx context.Context, pars *entities.OrderParameters) error {
	return u.srv.Delete(ctx, pars)
}

func (u *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return u.srv.Exists(ctx, orderNumber)
}
