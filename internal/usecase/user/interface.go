package user

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
)

type UserServiceI interface {
	IsLoginTaken(ctx context.Context, login string) (bool, error)
	IsValidPassword(password string, plainPassword string) bool
	HashPassword(password string) (string, error)
	GetBalanceWithWithdrawn(ctx context.Context, pars *model.GetPars) (*model.UserBalance, error)
	Register(ctx context.Context, user *model.GetPars) (*model.User, error)
	List(ctx context.Context, pars *model.ListPars) ([]*model.User, error)
	Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error)
	Create(ctx context.Context, obj *model.GetPars) error
	Update(ctx context.Context, pars *model.GetPars) error
	Delete(ctx context.Context, pars *model.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
