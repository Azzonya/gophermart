package user

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/user"
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

func (u *Usecase) Register(ctx context.Context, user *user.GetPars) (*user.User, error) {
	return u.srv.Register(ctx, user)
}

func (u *Usecase) IsValidPassword(password string, plainPassword string) bool {
	return u.srv.IsValidPassword(password, plainPassword)
}

func (u *Usecase) HashPassword(password string) (string, error) {
	return u.srv.HashPassword(password)
}

func (u *Usecase) CheckAuth(ctx context.Context, pars *user.GetPars) (*user.User, bool, error) {
	var result bool

	userFound, exist, err := u.srv.Get(ctx, &user.GetPars{
		Login: pars.Login,
	})
	if err != nil {
		return nil, false, err
	}
	if exist {
		result = u.IsValidPassword(userFound.Password, pars.Password)
	}

	return userFound, result, err
}

func (u *Usecase) GetBalanceWithWithdrawn(ctx context.Context, pars *user.GetPars) (*user.UserBalance, error) {
	return u.srv.GetBalanceWithWithdrawn(ctx, pars)
}

func (u *Usecase) List(ctx context.Context, pars *user.ListPars) ([]*user.User, error) {
	return u.srv.List(ctx, pars)
}

func (u *Usecase) Create(ctx context.Context, obj *user.GetPars) error {
	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Get(ctx context.Context, pars *user.GetPars) (*user.User, bool, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) GetById(ctx context.Context, userID string) (*user.User, bool, error) {
	return u.srv.Get(ctx, &user.GetPars{
		ID: userID,
	})
}

func (u *Usecase) Update(ctx context.Context, pars *user.GetPars) error {
	return u.srv.Update(ctx, pars)
}

func (u *Usecase) Delete(ctx context.Context, pars *user.GetPars) error {
	return u.srv.Delete(ctx, pars)
}

func (u *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return u.srv.Exists(ctx, orderNumber)
}
