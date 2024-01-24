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

func (u *Usecase) Register(ctx context.Context, user *model.GetPars) error {
	return u.srv.Register(ctx, user)
}

func (u *Usecase) CheckAuth(ctx context.Context, pars *model.GetPars) (bool, error) {
	_, exist, err := u.srv.Get(ctx, pars)
	return exist, err
}

func (s *Usecase) List(ctx context.Context, pars *model.ListPars) ([]*model.User, error) {
	return s.srv.List(ctx, pars)
}

func (s *Usecase) Create(ctx context.Context, obj *model.GetPars) error {
	return s.srv.Create(ctx, obj)
}

func (s *Usecase) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	return s.srv.Get(ctx, pars)
}

func (s *Usecase) Update(ctx context.Context, pars *model.GetPars, obj *model.GetPars) error {
	return s.srv.Update(ctx, pars, obj)
}

func (s *Usecase) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.srv.Delete(ctx, pars)
}

func (s *Usecase) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.srv.Exists(ctx, orderNumber)
}
