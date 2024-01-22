package user

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
)

type UserServiceI interface {
	IsLoginTaken(ctx context.Context, login string) (bool, error)
	Register(ctx context.Context, user *model.GetPars) error
	Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error)
}
