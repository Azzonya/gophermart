package bonustransactions

import (
	"context"
	"fmt"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BonusTransactionsRepoDBI interface {
	ListBt(ctx context.Context, pars *entities.BonusTransactionsListPars) ([]*entities.BonusTransaction, error)
	Create(ctx context.Context, obj *entities.BonusTransaction) error
	Get(ctx context.Context, pars *entities.BonusTransactionsParameters) (*entities.BonusTransaction, error)
	Update(ctx context.Context, pars *entities.BonusTransactionsParameters) error
	Delete(ctx context.Context, pars *entities.BonusTransactionsParameters) error
	Exists(ctx context.Context, login string) (bool, error)
}

type Repo struct {
	Con *pgxpool.Pool
}

func NewRepo(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

func (r *Repo) ListBt(ctx context.Context, pars *entities.BonusTransactionsListPars) ([]*entities.BonusTransaction, error) {
	queryBuilder := squirrel.Select("*").From("bonus_transactions").Where(squirrel.Eq{"true": true})

	if pars.ProcessedAfter != nil && !pars.ProcessedAfter.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"processed_at": pars.ProcessedAfter})
	}

	if pars.ProcessedBefore != nil && !pars.ProcessedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"processed_at": pars.ProcessedBefore})
	}

	if pars.TransactionType != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"transaction_type": pars.TransactionType})
	}

	if pars.MaxSum != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"sum": pars.MaxSum})
	}

	if pars.MinSum != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"sum": pars.MinSum})
	}

	if pars.UserID != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"user_id": pars.UserID})
	}

	if len(pars.OrderBy) != 0 {
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("processed_at %s", pars.OrderBy))
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Con.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*entities.BonusTransaction
	for rows.Next() {
		bonusTransaction := entities.BonusTransaction{}
		err = rows.Scan(&bonusTransaction.OrderNumber, &bonusTransaction.UserID, &bonusTransaction.ProcessedAt, &bonusTransaction.TransactionType, &bonusTransaction.Sum)
		if err != nil {
			return nil, err
		}

		result = append(result, &bonusTransaction)
	}

	return result, nil
}

func (r *Repo) Create(ctx context.Context, obj *entities.BonusTransaction) error {
	insert := squirrel.Insert("bonus_transactions").
		Columns("order_code", "user_id", "transaction_type", "sum").
		Values(obj.OrderNumber, obj.UserID, obj.TransactionType, obj.Sum)

	query, args, err := insert.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, query, args...)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *entities.BonusTransactionsParameters) (*entities.BonusTransaction, error) {
	queryBuilder := squirrel.Select("*").From("bonus_transactions").Where(squirrel.Eq{"true": true})

	if len(pars.OrderNumber) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_code": pars.OrderNumber})
	}

	if !pars.ProcessedAt.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"processed_at": pars.ProcessedAt})
	}

	if pars.Sum != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"sum": pars.Sum})
	}

	if len(pars.TransactionType) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"transaction_type": pars.TransactionType})
	}

	if len(pars.UserID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Con.QueryRow(ctx, sql, args...)

	var result entities.BonusTransaction
	err = row.Scan(&result.OrderNumber, &result.UserID, &result.ProcessedAt, &result.TransactionType, &result.Sum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *Repo) Update(ctx context.Context, pars *entities.BonusTransactionsParameters) error {
	queryBuilder := squirrel.Update("bonustransactions")

	queryBuilder = queryBuilder.Where(squirrel.Eq{"order_code": pars.OrderNumber})

	if !pars.ProcessedAt.IsZero() {
		queryBuilder = queryBuilder.Set("processed_at", pars.ProcessedAt)
	}

	if pars.Sum >= 0 {
		queryBuilder = queryBuilder.Set("sum", pars.Sum)
	}

	if len(pars.TransactionType) != 0 {
		queryBuilder = queryBuilder.Set("transaction_type", pars.TransactionType)
	}

	if len(pars.UserID) != 0 {
		queryBuilder = queryBuilder.Set("user_id", pars.UserID)
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

func (r *Repo) Delete(ctx context.Context, pars *entities.BonusTransactionsParameters) error {
	deleteQuery := squirrel.Delete("bonus_transactions").
		Where(squirrel.Eq{"order_code": pars.OrderNumber})

	query, args, err := deleteQuery.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, query, args...)
	return err
}

func (r *Repo) Exists(ctx context.Context, orderNumber string) (bool, error) {
	existsQuery := squirrel.Select("EXISTS (SELECT 1 FROM bonus_transactions WHERE order_code = ?)", orderNumber)

	query, args, err := existsQuery.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return false, err
	}

	var exist bool
	err = r.Con.QueryRow(ctx, query, args...).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, nil
}
