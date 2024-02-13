package bonustransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonustransactions"
)

type WithdrawalServiceI interface {
	List(ctx context.Context, pars *bonustransactions.ListPars) ([]*bonustransactions.BonusTransaction, error)
	Get(ctx context.Context, pars *bonustransactions.GetPars) (*bonustransactions.BonusTransaction, error)
	Create(ctx context.Context, obj *bonustransactions.BonusTransaction) error
	Update(ctx context.Context, pars *bonustransactions.GetPars) error
	Delete(ctx context.Context, pars *bonustransactions.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}

type BonusServiceI interface {
	WithdrawBalance(ctx context.Context, pars *bonustransactions.BonusTransaction) error
}
