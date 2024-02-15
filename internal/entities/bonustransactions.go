package entities

import (
	"time"
)

type BonusTransaction struct {
	OrderNumber     string
	UserID          string
	ProcessedAt     time.Time
	TransactionType TransactionType
	Sum             float32
}

type WithdrawalsResult struct {
	OrderNumber string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type TransactionType string

const (
	Accrual TransactionType = "+"
	Debit   TransactionType = "-"
)

type BonusTransactionsParameters struct {
	OrderNumber     string
	UserID          string
	ProcessedAt     time.Time
	TransactionType TransactionType
	Sum             float32
}

type BonusTransactionsListPars struct {
	OrderNumber     *string
	UserID          *string
	ProcessedBefore *time.Time
	ProcessedAfter  *time.Time
	TransactionType TransactionType
	MinSum          *float32
	MaxSum          *float32
	OrderBy         string
}
