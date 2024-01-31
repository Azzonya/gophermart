package bonus_transactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonus_transactions/model"
)

type Usecase struct {
	srv WithdrawalServiceI
}

func New(srv WithdrawalServiceI) *Usecase {
	return &Usecase{
		srv: srv,
	}
}

func (u *Usecase) List(ctx context.Context, pars *model.ListPars) ([]*model.BonusTransaction, error) {
	return u.srv.List(ctx, pars)
}

func (u *Usecase) Create(ctx context.Context, obj *model.GetPars) error {
	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Get(ctx context.Context, pars *model.GetPars) (*model.BonusTransaction, bool, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) Update(ctx context.Context, pars *model.GetPars) error {
	return u.srv.Update(ctx, pars)
}

func (u *Usecase) Delete(ctx context.Context, pars *model.GetPars) error {
	return u.srv.Delete(ctx, pars)
}

func (u *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return u.srv.Exists(ctx, orderNumber)
}
