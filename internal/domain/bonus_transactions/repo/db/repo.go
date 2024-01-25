package db

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/bonus_transactions/model"
	commonRepoPg "github.com/Azzonya/gophermart/internal/domain/common/repo/pg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	*commonRepoPg.Base
}

func New(con *pgxpool.Pool) *Repo {
	return &Repo{
		commonRepoPg.NewBase(con),
	}
}

func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.BonusTransaction, error) {
	var result []*model.BonusTransaction
	var values []interface{}
	query := "SELECT * FROM bonus_transactions WHERE true"

	if !pars.ProcessedAfter.IsZero() {
		query += " AND processed_at >= $1"
		values = append(values, pars.ProcessedAfter)
	}

	if !pars.ProcessedBefore.IsZero() {
		query += " AND processed_at <= $2"
		values = append(values, pars.ProcessedBefore)
	}

	if pars.TransactionType != nil {
		query += " AND transaction_type = $3"
		values = append(values, *pars.TransactionType)
	}

	if pars.MaxSum != nil {
		query += " AND sum <= $4"
		values = append(values, *pars.MaxSum)
	}

	if pars.MinSum != nil {
		query += " AND sum >= $5"
		values = append(values, *pars.MinSum)
	}

	rows, err := r.Con.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var bonusTransaction *model.BonusTransaction
		err = rows.Scan(&bonusTransaction.OrderNumber, &bonusTransaction.UserID, &bonusTransaction.ProcessedAt, &bonusTransaction.TransactionType, &bonusTransaction.Sum)
		if err != nil {
			return nil, err
		}

		result = append(result, bonusTransaction)
	}

	return result, err
}

func (r *Repo) Create(ctx context.Context, obj *model.GetPars) error {
	_, err := r.Con.Exec(ctx, "INSERT INTO bonus_transactions (order_code, user_id, transacton_type, sum) VALUES ($1, $2, $3, $4);", obj.OrderNumber, obj.UserID, obj.TransactionType, obj.Sum)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.BonusTransaction, bool, error) {
	var values []interface{}
	var result model.BonusTransaction
	query := "SELECT * FROM bonus_transactions WHERE true"

	if len(pars.OrderNumber) != 0 {
		query += " AND sum = $1"
		values = append(values, pars.Sum)
	}

	if !pars.ProcessedAt.IsZero() {
		query += " AND processed_at = $2"
		values = append(values, pars.ProcessedAt)
	}

	if pars.Sum != 0 {
		query += " AND sum = $3"
		values = append(values, pars.Sum)
	}

	if len(pars.TransactionType) != 0 {
		query += " AND transaction_type = $4"
		values = append(values, pars.TransactionType)
	}

	if pars.UserID != 0 {
		query += " AND user_id = $5"
		values = append(values, pars.UserID)
	}

	err := r.Con.QueryRow(ctx, query, values...).Scan(&result)
	if err != nil {
		return nil, false, err
	}

	return &result, result.OrderNumber != "", err
}

func (r *Repo) Update(ctx context.Context, pars *model.GetPars) error {
	var values []interface{}

	query := "UPDATE bonus_transactions SET 1=1"

	if !pars.ProcessedAt.IsZero() {
		query += " AND processed_at = $1"
		values = append(values, pars.ProcessedAt)
	}

	if pars.Sum != 0 {
		query += " AND sum = $2"
		values = append(values, pars.Sum)
	}

	if len(pars.TransactionType) != 0 {
		query += " AND transaction_type = $3"
		values = append(values, pars.TransactionType)
	}

	if pars.UserID != 0 {
		query += " AND user_id = $4"
		values = append(values, pars.UserID)
	}

	_, err := r.Con.Exec(ctx, query, values...)

	return err
}

func (r *Repo) Delete(ctx context.Context, pars *model.GetPars) error {
	_, err := r.Con.Exec(ctx, "DELETE FROM bonus_transactions WHERE order_code = $1;", pars.OrderNumber)
	return err
}

func (r *Repo) Exists(ctx context.Context, orderNumber string) (bool, error) {
	var exist bool
	err := r.Con.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM bonus_transactions WHERE order_code = $1);", orderNumber).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, err
}
