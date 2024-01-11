package model

import "time"

type Order struct {
	OrderNumber int
	UploadedAt  time.Time
	Status      string
	UserID      string
	Accrual     int
}
