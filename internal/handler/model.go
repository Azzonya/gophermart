package handler

import "github.com/Azzonya/gophermart/internal/domain/user/model"

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Decode() *model.User {
	return &model.User{Login: r.Login, Password: r.Password}
}
