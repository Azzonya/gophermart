package handler

import (
	orderModel "github.com/Azzonya/gophermart/internal/domain/order/model"
	userModel "github.com/Azzonya/gophermart/internal/domain/user/model"
	"time"
)

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Decode() *userModel.User {
	return &userModel.User{Login: r.Login, Password: r.Password}
}

type ListOrdersResult struct {
	OrderNumber string                 `json:"number"`
	Status      orderModel.OrderStatus `json:"status"`
	Accrual     int                    `json:"accrual,omitempty"`
	UploadedAt  time.Time              `json:"uploaded_at"`
}

func (l *ListOrdersResult) Encode(order *orderModel.Order) {
	//l.UploadedAt = order.UploadedAt
	//l.Status = order.Status
	//l.Accrual = order.
	//l.OrderNumber = order.OrderNumber
}

type UserBalanceResult struct {
	Balance   float64 `json:"balance"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawalBalanceRequest struct {
	OrderNumber string `json:"order"`
	Sum         int    `json:"sum"`
}
