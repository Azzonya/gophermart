package user

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
)

type Usecase struct {
	srv UserServiceI
}

func New(srv UserServiceI) *Usecase {
	return &Usecase{
		srv: srv,
	}
}

func (u *Usecase) IsLoginTaken(ctx context.Context, login string) (bool, error) {
	return u.srv.IsLoginTaken(ctx, login)
}

func (u *Usecase) Register(ctx context.Context, user *model.User) error {
	return nil
}
