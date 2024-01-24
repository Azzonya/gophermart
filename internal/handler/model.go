package handler

import (
	model2 "github.com/Azzonya/gophermart/internal/domain/order/model"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
	"time"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Decode() *model.User {
	return &model.User{Login: r.Login, Password: r.Password}
}

type ListOrdersResult struct {
	OrderNumber string    `json:"number"`
	Status      string    `json:"status"`
	Accrual     int       `json:"accrual,omitempty"`
	UploadedAt  time.Time `json:"uploaded_at"`
}

func (l *ListOrdersResult) Encode(order *model2.Order) {
	l.UploadedAt = order.UploadedAt
	l.Status = order.Status
	l.Accrual = order.Accrual
	l.OrderNumber = order.OrderNumber
}

type UserBalanceResult struct {
	Balance   float64 `json:"balance"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalBalanceRequest struct {
	OrderNumber string `json:"order"`
	Sum         int    `json:"sum"`
}
