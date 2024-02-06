package bonusTransactions

import (
	"context"
	"time"
)

type RepoDBIntreface interface {
	List(ctx context.Context, pars *ListPars) ([]*BonusTransaction, error)
	Create(ctx context.Context, obj *GetPars) error
	Get(ctx context.Context, pars *GetPars) (*BonusTransaction, bool, error)
	Update(ctx context.Context, pars *GetPars) error
	Delete(ctx context.Context, pars *GetPars) error
	Exists(ctx context.Context, login string) (bool, error)
}

type BonusTransaction struct {
	OrderNumber     string
	UserID          string
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
	UserID          string
	ProcessedAt     time.Time
	TransactionType TransactionType
	Sum             int
}

type ListPars struct {
	OrderNumber     *string
	UserID          *string
	ProcessedBefore *time.Time
	ProcessedAfter  *time.Time
	TransactionType TransactionType
	MinSum          *int
	MaxSum          *int
	OrderBy         string
}
