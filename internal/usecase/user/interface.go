package user

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user"
)

type UserServiceI interface {
	IsLoginTaken(ctx context.Context, login string) (bool, error)
	IsValidPassword(password string, plainPassword string) bool
	HashPassword(password string) (string, error)
	GetBalanceWithWithdrawn(ctx context.Context, pars *user.GetPars) (*user.UserBalance, error)
	Register(ctx context.Context, user *user.GetPars) (*user.User, error)
	List(ctx context.Context, pars *user.ListPars) ([]*user.User, error)
	Get(ctx context.Context, pars *user.GetPars) (*user.User, error)
	Create(ctx context.Context, obj *user.GetPars) error
	Update(ctx context.Context, pars *user.GetPars) error
	Delete(ctx context.Context, pars *user.GetPars) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
