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

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
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
	UserID         *string
}
