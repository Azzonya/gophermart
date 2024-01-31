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

func (u *Usecase) Register(ctx context.Context, user *model.GetPars) (*model.User, error) {
	return u.srv.Register(ctx, user)
}

func (u *Usecase) IsValidPassword(password string, plainPassword string) bool {
	return u.srv.IsValidPassword(password, plainPassword)
}

func (u *Usecase) HashPassword(password string) (string, error) {
	return u.srv.HashPassword(password)
}

func (u *Usecase) CheckAuth(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	var result bool

	user, exist, err := u.srv.Get(ctx, &model.GetPars{
		Login: pars.Login,
	})
	if err != nil {
		return nil, false, err
	}
	if exist {
		result = u.IsValidPassword(user.Password, pars.Password)
	}

	return user, result, err
}

func (u *Usecase) GetBalanceWithWithdrawn(ctx context.Context, pars *model.GetPars) (*model.UserBalance, error) {
	return u.srv.GetBalanceWithWithdrawn(ctx, pars)
}

func (u *Usecase) List(ctx context.Context, pars *model.ListPars) ([]*model.User, error) {
	return u.srv.List(ctx, pars)
}

func (u *Usecase) Create(ctx context.Context, obj *model.GetPars) error {
	return u.srv.Create(ctx, obj)
}

func (u *Usecase) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	return u.srv.Get(ctx, pars)
}

func (u *Usecase) Update(ctx context.Context, pars *model.GetPars) error {
	return u.srv.Update(ctx, pars)
}

func (u *Usecase) Delete(ctx context.Context, pars *model.GetPars) error {
	return u.srv.Delete(ctx, pars)
}

func (u *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return u.srv.Exists(ctx, orderNumber)
}
