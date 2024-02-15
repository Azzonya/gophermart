package bonustransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
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

func (u *Usecase) List(ctx context.Context, pars *entities.BonusTransactionsListPars) ([]*entities.WithdrawalsResult, error) {
	withdrawals, err := u.srv.ListBtS(ctx, pars)
	if err != nil {
		return nil, err
	}

	result := []*entities.WithdrawalsResult{}
	for _, v := range withdrawals {
		withdrawal := entities.WithdrawalsResult{
			OrderNumber: v.OrderNumber,
			Sum:         v.Sum,
			ProcessedAt: v.ProcessedAt.Format(time.RFC3339),
		}

		result = append(result, &withdrawal)
	}

	return result, nil
}

func (u *Usecase) Create(ctx context.Context, obj *entities.BonusTransaction) error {
	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Get(ctx context.Context, pars *entities.BonusTransactionsParameters) (*entities.BonusTransaction, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) Update(ctx context.Context, pars *entities.BonusTransactionsParameters) error {
	return u.srv.Update(ctx, pars)
}

func (u *Usecase) Delete(ctx context.Context, pars *entities.BonusTransactionsParameters) error {
	return u.srv.Delete(ctx, pars)
}

func (u *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return u.srv.Exists(ctx, orderNumber)
}

func (u *Usecase) WithdrawBalanceU(ctx context.Context, pars *entities.BonusTransaction) error {
	return u.bonus.WithdrawBalance(ctx, pars)
}
