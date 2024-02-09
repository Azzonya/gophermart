package handler

import (
	userModel "github.com/Azzonya/gophermart/internal/domain/user"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Decode() *userModel.User {
	return &userModel.User{Login: r.Login, Password: r.Password}
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

type WithdrawalsResult struct {
	OrderNumber string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
