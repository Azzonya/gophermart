package user

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
)

type UserServiceI interface {
	IsLoginTaken(ctx context.Context, login string) (bool, error)
	IsValidPassword(password string, plainPassword string) bool
	HashPassword(password string) (string, error)
	GetBalanceWithWithdrawn(ctx context.Context, pars *entities.UserParameters) (*entities.UserBalance, error)
	Register(ctx context.Context, user *entities.User) (*entities.User, error)
	List(ctx context.Context, pars *entities.UserListPars) ([]*entities.User, error)
	Get(ctx context.Context, pars *entities.UserParameters) (*entities.User, error)
	Create(ctx context.Context, obj *entities.User) error
	Update(ctx context.Context, pars *entities.UserParameters) error
	Delete(ctx context.Context, pars *entities.UserParameters) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}
