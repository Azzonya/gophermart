package bonustransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonustransactions"
	"time"
)

type Usecase struct {
	srv   WithdrawalServiceI
	bonus BonusServiceI
}

func New(srv WithdrawalServiceI) *Usecase {
	return &Usecase{
		srv: srv,
	}
}

func (u *Usecase) SetBonusService(bonus BonusServiceI) {
	u.bonus = bonus
}

func (u *Usecase) List(ctx context.Context, pars *bonustransactions.ListPars) ([]*bonustransactions.WithdrawalsResult, error) {
	withdrawals, err := u.srv.List(ctx, pars)
	if err != nil {
		return nil, err
	}

	result := []*bonustransactions.WithdrawalsResult{}
	for _, v := range withdrawals {
		withdrawal := bonustransactions.WithdrawalsResult{
			OrderNumber: v.OrderNumber,
			Sum:         v.Sum,
			ProcessedAt: v.ProcessedAt.Format(time.RFC3339),
		}

		result = append(result, &withdrawal)
	}

	return result, nil
}

func (u *Usecase) Create(ctx context.Context, obj *bonustransactions.BonusTransaction) error {
	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Get(ctx context.Context, pars *bonustransactions.GetPars) (*bonustransactions.BonusTransaction, error) {
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

func (u *Usecase) WithdrawBalance(ctx context.Context, pars *bonustransactions.BonusTransaction) error {
	return u.bonus.WithdrawBalance(ctx, pars)
}
