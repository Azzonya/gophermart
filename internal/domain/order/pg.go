package order

import (
	"context"
	"errors"
	"fmt"
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

func (r *Repo) List(ctx context.Context, pars *ListPars) ([]*Order, error) {
	var result []*Order
	var values []interface{}
	values = make([]interface{}, 0)

	query := "SELECT code, uploaded_at, status, user_id FROM orders WHERE true"

	paramNum := 1
	if pars.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", paramNum)
		values = append(values, *pars.UserID)
		paramNum++
	}

	if pars.OrderNumber != nil {
		query += fmt.Sprintf(" AND code = $%d OR $%d IS NULL", paramNum, paramNum)
		values = append(values, *pars.OrderNumber)
		paramNum++
	}

	if pars.UploadedBefore != nil {
		query += fmt.Sprintf(" AND uploaded_at <= $%d", paramNum)
		values = append(values, *pars.UploadedBefore)
		paramNum++
	}

	if pars.UploadedAfter != nil {
		query += fmt.Sprintf(" AND uploaded_at >= $%d", paramNum)
		values = append(values, *pars.UploadedAfter)
		paramNum++
	}

	if pars.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", paramNum)
		values = append(values, *pars.Status)
		paramNum++
	}

	if pars.Statuses != nil && len(pars.Statuses) > 0 {
		query += " AND status IN ("
		for i, v := range pars.Statuses {
			query += fmt.Sprintf("'%s'", v)
			if i < len(pars.Statuses)-1 {
				query += ","
			}
		}
		query += ")"
	}

	if len(pars.OrderBy) != 0 {
		query += fmt.Sprintf(" ORDER BY uploaded_at %s", pars.OrderBy)
	}

	rows, err := r.Con.Query(ctx, query, values...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var orderFound Order
		err = rows.Scan(&orderFound.OrderNumber, &orderFound.UploadedAt, &orderFound.Status, &orderFound.UserID)
		if err != nil {
			return nil, err
		}

		result = append(result, &orderFound)
	}

	return result, err
}

func (r *Repo) Create(ctx context.Context, obj *Order) error {
	_, err := r.Con.Exec(ctx, "INSERT INTO orders (code, status, user_id) VALUES ($1, $2, $3);", obj.OrderNumber, obj.Status, obj.UserID)

	return err
}

func (r *Repo) Get(ctx context.Context, pars *GetPars) (*Order, error) {
	var values []interface{}
	var result Order
	values = make([]interface{}, 0)

	query := "SELECT * FROM orders WHERE true"

	paramNum := 1
	if len(pars.UserID) != 0 {
		query += fmt.Sprintf(" AND user_id = $%d", paramNum)
		values = append(values, pars.UserID)
		paramNum += 1
	}

	if len(pars.OrderNumber) != 0 {
		query += fmt.Sprintf(" AND code = $%d", paramNum)
		values = append(values, pars.OrderNumber)
		paramNum += 1
	}

	if len(pars.Status) != 0 {
		query += fmt.Sprintf(" AND status = $%d", paramNum)
		values = append(values, pars.Status)
		paramNum += 1
	}

	if len(pars.UserID) != 0 {
		query += fmt.Sprintf(" AND user_id >= $%d", paramNum)
		values = append(values, pars.UserID)
		paramNum += 1
	}
	err := r.Con.QueryRow(ctx, query, values...).Scan(&result.OrderNumber, &result.UploadedAt, &result.Status, &result.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *Repo) Update(ctx context.Context, pars *GetPars) error {
	var values []interface{}
	values = make([]interface{}, 0)

	query := "UPDATE orders"

	paramNum := 1
	if len(pars.UserID) != 0 {
		query += fmt.Sprintf(" SET user_id = $%d", paramNum)
		values = append(values, pars.UserID)
		paramNum++
	}

	if len(pars.Status) != 0 {
		if len(values) > 0 {
			query += ","
		} else {
			query += " SET"
		}
		query += fmt.Sprintf(" status = $%d", paramNum)
		values = append(values, pars.Status)
		paramNum++
	}

	query += fmt.Sprintf(" WHERE code = '%s'", pars.OrderNumber)

	_, err := r.Con.Exec(ctx, query, values...)

	return err
}

func (r *Repo) Delete(ctx context.Context, pars *GetPars) error {
	_, err := r.Con.Exec(ctx, "DELETE FROM orders WHERE code = $1;", pars.OrderNumber)
	return err
}

func (r *Repo) Exists(ctx context.Context, orderNumber string) (bool, error) {
	var exist bool
	err := r.Con.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM orders WHERE code = $1);", orderNumber).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, err
}
