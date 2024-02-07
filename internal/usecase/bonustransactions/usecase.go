package bonustransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonustransactions"
)

type Usecase struct {
	srv WithdrawalServiceI
}

func New(srv WithdrawalServiceI) *Usecase {
	return &Usecase{
		srv: srv,
	}
}

func (u *Usecase) List(ctx context.Context, pars *bonustransactions.ListPars) ([]*bonustransactions.BonusTransaction, error) {
	return u.srv.List(ctx, pars)
}

func (u *Usecase) Create(ctx context.Context, obj *bonustransactions.GetPars) error {
	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Get(ctx context.Context, pars *bonustransactions.GetPars) (*bonustransactions.BonusTransaction, bool, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) Update(ctx context.Context, pars *bonustransactions.GetPars) error {
	return u.srv.Update(ctx, pars)
}

func (u *Usecase) Delete(ctx context.Context, pars *bonustransactions.GetPars) error {
	return u.srv.Delete(ctx, pars)
}

func (u *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return u.srv.Exists(ctx, orderNumber)
}
