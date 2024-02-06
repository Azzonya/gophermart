package bonustransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonusTransactions"
)

type WithdrawalServiceI interface {
	List(ctx context.Context, pars *bonusTransactions.ListPars) ([]*bonusTransactions.BonusTransaction, error)
	Get(ctx context.Context, pars *bonusTransactions.GetPars) (*bonusTransactions.BonusTransaction, bool, error)
	Create(ctx context.Context, obj *bonusTransactions.GetPars) error
	Update(ctx context.Context, pars *bonusTransactions.GetPars) error
	Delete(ctx context.Context, pars *bonusTransactions.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
