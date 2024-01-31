package model

import (
	"time"
)

type Order struct {
	OrderNumber string
	UploadedAt  time.Time
	Status      OrderStatus
	UserID      string
}

type OrderWithAccrual struct {
	OrderNumber string      `json:"number"`
	Status      OrderStatus `json:"status"`
	Accrual     int         `json:"accrual,omitempty"`
	UploadedAt  string      `json:"uploaded_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	//OrderStatusInvalid    OrderStatus = "INVALID"
	//OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type GetPars struct {
	OrderNumber string
	UploadedAt  time.Time
	Status      OrderStatus
	UserID      string
}

type ListPars struct {
	OrderNumber    *string
	UploadedBefore *time.Time
	UploadedAfter  *time.Time
	Status         *OrderStatus
	Statuses       []OrderStatus
	UserID         *string
	OrderBy        string
}
