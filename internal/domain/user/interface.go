package user

import (
	"context"
)

type repoDBI interface {
	List(ctx context.Context, pars *ListPars) ([]*User, error)
	Create(ctx context.Context, obj *GetPars) error
	Get(ctx context.Context, pars *GetPars) (*User, bool, error)
	Update(ctx context.Context, pars *GetPars) error
	Delete(ctx context.Context, pars *GetPars) error
	Exists(ctx context.Context, login string) (bool, error)
}
