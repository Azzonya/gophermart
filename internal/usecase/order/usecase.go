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

func (u *Usecase) Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error) {

	return nil, false, nil
}
