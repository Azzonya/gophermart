package order

import (
	"context"
)

type RepoDBI interface {
	List(ctx context.Context, pars *ListPars) ([]*Order, error)
	Create(ctx context.Context, obj *GetPars) error
	Get(ctx context.Context, pars *GetPars) (*Order, bool, error)
	Update(ctx context.Context, pars *GetPars) error
	Delete(ctx context.Context, pars *GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
