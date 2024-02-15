package bonustransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
)

type WithdrawalServiceI interface {
	List(ctx context.Context, pars *entities.BonusTransactionsListPars) ([]*entities.BonusTransaction, error)
	Get(ctx context.Context, pars *entities.BonusTransactionsParameters) (*entities.BonusTransaction, error)
	Create(ctx context.Context, obj *entities.BonusTransaction) error
	Update(ctx context.Context, pars *entities.BonusTransactionsParameters) error
	Delete(ctx context.Context, pars *entities.BonusTransactionsParameters) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}

type BonusServiceI interface {
	WithdrawBalance(ctx context.Context, pars *entities.BonusTransaction) error
}
