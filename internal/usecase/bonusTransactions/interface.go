package bonusTransactions

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonusTransactions/model"
)

type WithdrawalServiceI interface {
	List(ctx context.Context, pars *model.ListPars) ([]*model.BonusTransaction, error)
	Get(ctx context.Context, pars *model.GetPars) (*model.BonusTransaction, bool, error)
	Create(ctx context.Context, obj *model.GetPars) error
	Update(ctx context.Context, pars *model.GetPars) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
