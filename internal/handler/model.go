package handler

import (
	"github.com/Azzonya/gophermart/internal/entities"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Decode() *entities.User {
	return &entities.User{Login: r.Login, Password: r.Password}
}

type UserBalanceResult struct {
	Balance   float32 `json:"balance"`
	Withdrawn float32 `json:"withdrawn"`
}

type WithdrawalBalanceRequest struct {
	UserID      string
	OrderNumber string  `json:"order"`
	Sum         float32 `json:"sum"`
}
