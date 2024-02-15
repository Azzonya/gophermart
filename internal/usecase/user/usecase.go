package user

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Usecase struct {
	srv     UserServiceI
	authSrv AuthServiceI
}

func New(srv UserServiceI, authServ AuthServiceI) *Usecase {
	return &Usecase{
		srv:     srv,
		authSrv: authServ,
	}
}

func (u *Usecase) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	return u.srv.Register(ctx, user)
}

func (u *Usecase) IsValidPassword(password string, plainPassword string) bool {
	return u.srv.IsValidPassword(password, plainPassword)
}

func (u *Usecase) CheckAuth(ctx context.Context, pars *entities.UserParameters) (*entities.User, bool, error) {
	var result bool

	userFound, err := u.srv.Get(ctx, &entities.UserParameters{
		Login: pars.Login,
	})
	if err != nil {
		return nil, false, err
	}
	if userFound != nil {
		result = u.IsValidPassword(userFound.Password, pars.Password)
	}

	return userFound, result, err
}

func (u *Usecase) GetBalanceWithWithdrawn(ctx context.Context, pars *entities.UserParameters) (*entities.UserBalance, error) {
	return u.srv.GetBalanceWithWithdrawn(ctx, pars)
}

func (u *Usecase) GetUserIDFromCookieU(ctx *gin.Context) (string, error) {
	return u.authSrv.GetUserIDFromCookie(ctx)
}

func (u *Usecase) CreateJWTCookieU(user *entities.User) (*http.Cookie, error) {
	return u.authSrv.CreateJWTCookie(user)
}
