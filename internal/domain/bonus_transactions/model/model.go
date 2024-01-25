package model

import "time"

type BonusTransaction struct {
	OrderNumber     string
	UserID          int
	ProcessedAt     time.Time
	TransactionType TransactionType
	Sum             int
}

type TransactionType string

const (
	Accrual TransactionType = "+"
	Debit   TransactionType = "-"
)

type GetPars struct {
	OrderNumber     string
	UserID          int
	ProcessedAt     time.Time
	TransactionType TransactionType
	Sum             int
}

type ListPars struct {
	OrderNumber     *string
	UserID          *int
	ProcessedBefore *time.Time
	ProcessedAfter  *time.Time
	TransactionType *TransactionType
	MinSum          *int
	MaxSum          *int
}
