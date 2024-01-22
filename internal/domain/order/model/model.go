package model

import (
	"time"
)

type Order struct {
	OrderNumber string
	UploadedAt  time.Time
	Status      string
	UserID      string
	Accrual     int
}

type GetPars struct {
	OrderNumber string
	UploadedAt  time.Time
	Status      string
	UserID      string
	Accrual     int
}

type ListPars struct {
	OrderNumber    *string
	UploadedBefore *time.Time
	UploadedAfter  *time.Time
	Status         *string
	UserID         *string
	MinAccrual     *int
	MaxAccrual     *int
}
