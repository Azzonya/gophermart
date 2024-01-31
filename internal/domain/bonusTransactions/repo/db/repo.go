package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azzonya/gophermart/internal/domain/bonusTransactions/model"
	commonRepoPg "github.com/Azzonya/gophermart/internal/domain/common/repo/pg"
	"github.com/jackc/pgx/v5"
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

	paramNum := 1
	if pars.ProcessedAfter != nil && !pars.ProcessedAfter.IsZero() {
		query += fmt.Sprintf(" AND processed_at >= $%d", paramNum)
		values = append(values, pars.ProcessedAfter)
		paramNum += 1
	}

	if pars.ProcessedBefore != nil && !pars.ProcessedBefore.IsZero() {
		query += fmt.Sprintf(" AND processed_at <= $%d", paramNum)
		values = append(values, pars.ProcessedBefore)
		paramNum += 1
	}

	if pars.TransactionType != "" {
		query += fmt.Sprintf(" AND transaction_type = $%d", paramNum)
		values = append(values, pars.TransactionType)
		paramNum += 1
	}

	if pars.MaxSum != nil {
		query += fmt.Sprintf(" AND sum <= $%d", paramNum)
		values = append(values, *pars.MaxSum)
		paramNum += 1
	}

	if pars.MinSum != nil {
		query += fmt.Sprintf(" AND sum >= $%d", paramNum)
		values = append(values, *pars.MinSum)
		paramNum += 1
	}

	if pars.UserID != nil {
		query += fmt.Sprintf(" AND user_id >= $%d", paramNum)
		values = append(values, *pars.UserID)
		paramNum += 1
	}

	if len(pars.OrderBy) != 0 {
		query += fmt.Sprintf(" ORDER BY processed_at %s", pars.OrderBy)
	}

	rows, err := r.Con.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		bonusTransaction := model.BonusTransaction{}
		err = rows.Scan(&bonusTransaction.OrderNumber, &bonusTransaction.UserID, &bonusTransaction.ProcessedAt, &bonusTransaction.TransactionType, &bonusTransaction.Sum)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, nil
			}
			return nil, err
		}

		result = append(result, &bonusTransaction)
	}

	return result, err
}

func (r *Repo) Create(ctx context.Context, obj *model.GetPars) error {
	_, err := r.Con.Exec(ctx, "INSERT INTO bonus_transactions (order_code, user_id, transaction_type, sum) VALUES ($1, $2, $3, $4);", obj.OrderNumber, obj.UserID, obj.TransactionType, obj.Sum)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.BonusTransaction, bool, error) {
	var values []interface{}
	var result model.BonusTransaction
	query := "SELECT * FROM bonus_transactions WHERE true"

	paramNum := 1
	if len(pars.OrderNumber) != 0 {
		query += fmt.Sprintf(" AND sum = $%d", paramNum)
		values = append(values, pars.Sum)
		paramNum += 1
	}

	if !pars.ProcessedAt.IsZero() {
		query += fmt.Sprintf(" AND processed_at = $%d", paramNum)
		values = append(values, pars.ProcessedAt)
		paramNum += 1
	}

	if pars.Sum != 0 {
		query += fmt.Sprintf(" AND order_code = $%d", paramNum)
		values = append(values, pars.Sum)
		paramNum += 1
	}

	if len(pars.TransactionType) != 0 {
		query += fmt.Sprintf(" AND transaction_type = $%d", paramNum)
		values = append(values, pars.TransactionType)
		paramNum += 1
	}

	if len(pars.UserID) != 0 {
		query += fmt.Sprintf(" AND user_id = $%d", paramNum)
		values = append(values, pars.UserID)
		paramNum += 1
	}

	err := r.Con.QueryRow(ctx, query, values...).Scan(&result.OrderNumber, &result.UserID, &result.ProcessedAt, &result.TransactionType, &result.Sum)
	if err != nil {
		return nil, false, err
	}

	return &result, result.OrderNumber != "", err
}

func (r *Repo) Update(ctx context.Context, pars *model.GetPars) error {
	var values []interface{}

	query := "UPDATE bonustransactions"

	paramNum := 1
	if !pars.ProcessedAt.IsZero() {
		query += fmt.Sprintf(" SET processed_at = $%d", paramNum)
		values = append(values, pars.ProcessedAt)
		paramNum++
	}

	if pars.Sum >= 0 {
		if len(values) > 0 {
			query += ","
		} else {
			query += " SET"
		}
		query += fmt.Sprintf(" SET sum = $%d", paramNum)
		values = append(values, pars.Sum)
		paramNum++
	}

	if len(pars.TransactionType) != 0 {
		if len(values) > 0 {
			query += ","
		} else {
			query += " SET"
		}
		query += fmt.Sprintf(" SET transaction_type = $%d", paramNum)
		values = append(values, pars.TransactionType)
		paramNum++
	}

	if len(pars.UserID) != 0 {
		if len(values) > 0 {
			query += ","
		} else {
			query += " SET"
		}
		query += fmt.Sprintf(" SET user_id = $%d", paramNum)
		values = append(values, pars.UserID)
		paramNum++
	}

	query += fmt.Sprintf(" WHERE order_code = '%s'", pars.OrderNumber)

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
