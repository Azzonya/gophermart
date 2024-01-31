package service

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
)

type RepoDbI interface {
	List(ctx context.Context, pars *model.ListPars) ([]*model.User, error)
	Create(ctx context.Context, obj *model.GetPars) error
	Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error)
	Update(ctx context.Context, pars *model.GetPars) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, login string) (bool, error)
}
