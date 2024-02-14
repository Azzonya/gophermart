package bonustransactions

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Con *pgxpool.Pool
}

func NewRepo(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

func (r *Repo) List(ctx context.Context, pars *ListPars) ([]*BonusTransaction, error) {
	queryBuilder := squirrel.Select("*").From("bonus_transactions").Where(squirrel.Eq{"true": true})
	var values = make(map[string]interface{})

	if pars.ProcessedAfter != nil && !pars.ProcessedAfter.IsZero() {
		values["processed_at"] = pars.ProcessedAfter
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"processed_at": values["processed_at"]})
	}

	if pars.ProcessedBefore != nil && !pars.ProcessedBefore.IsZero() {
		values["processed_at"] = pars.ProcessedBefore
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"processed_at": values["processed_at"]})
	}

	if pars.TransactionType != "" {
		values["transaction_type"] = pars.TransactionType
		queryBuilder = queryBuilder.Where(squirrel.Eq{"transaction_type": values["transaction_type"]})
	}

	if pars.MaxSum != nil {
		values["sum"] = *pars.MaxSum
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"sum": values["sum"]})
	}

	if pars.MinSum != nil {
		values["sum"] = *pars.MinSum
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"sum": values["sum"]})
	}

	if pars.UserID != nil {
		values["user_id"] = *pars.UserID
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"user_id": values["user_id"]})
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

	var result []*BonusTransaction
	for rows.Next() {
		bonusTransaction := BonusTransaction{}
		err = rows.Scan(&bonusTransaction.OrderNumber, &bonusTransaction.UserID, &bonusTransaction.ProcessedAt, &bonusTransaction.TransactionType, &bonusTransaction.Sum)
		if err != nil {
			return nil, err
		}

		result = append(result, &bonusTransaction)
	}

	return result, nil
}

func (r *Repo) Create(ctx context.Context, obj *BonusTransaction) error {
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

func (r *Repo) Get(ctx context.Context, pars *GetPars) (*BonusTransaction, error) {
	queryBuilder := squirrel.Select("*").From("bonus_transactions").Where(squirrel.Eq{"true": true})

	var values = make(map[string]interface{})

	if len(pars.OrderNumber) != 0 {
		values["order_code"] = pars.OrderNumber
		queryBuilder = queryBuilder.Where(squirrel.Eq{"order_code": values["order_code"]})
	}

	if !pars.ProcessedAt.IsZero() {
		values["processed_at"] = pars.ProcessedAt
		queryBuilder = queryBuilder.Where(squirrel.Eq{"processed_at": values["processed_at"]})
	}

	if pars.Sum != 0 {
		values["sum"] = pars.Sum
		queryBuilder = queryBuilder.Where(squirrel.Eq{"sum": values["sum"]})
	}

	if len(pars.TransactionType) != 0 {
		values["transaction_type"] = pars.TransactionType
		queryBuilder = queryBuilder.Where(squirrel.Eq{"transaction_type": values["transaction_type"]})
	}

	if len(pars.UserID) != 0 {
		values["user_id"] = pars.UserID
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": values["user_id"]})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := r.Con.QueryRow(ctx, sql, args...)

	var result BonusTransaction
	err = row.Scan(&result.OrderNumber, &result.UserID, &result.ProcessedAt, &result.TransactionType, &result.Sum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *Repo) Update(ctx context.Context, pars *GetPars) error {
	queryBuilder := squirrel.Update("bonustransactions")

	queryBuilder = queryBuilder.Where(squirrel.Eq{"order_code": pars.OrderNumber})

	var values = make(map[string]interface{})

	if !pars.ProcessedAt.IsZero() {
		values["processed_at"] = pars.ProcessedAt
		queryBuilder = queryBuilder.Set("processed_at", values["processed_at"])
	}

	if pars.Sum >= 0 {
		values["sum"] = pars.Sum
		queryBuilder = queryBuilder.Set("sum", values["sum"])
	}

	if len(pars.TransactionType) != 0 {
		values["transaction_type"] = pars.TransactionType
		queryBuilder = queryBuilder.Set("transaction_type", values["transaction_type"])
	}

	if len(pars.UserID) != 0 {
		values["user_id"] = pars.UserID
		queryBuilder = queryBuilder.Set("user_id", values["user_id"])
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

func (r *Repo) Delete(ctx context.Context, pars *GetPars) error {
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
