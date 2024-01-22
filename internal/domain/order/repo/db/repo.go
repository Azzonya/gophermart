package db

import (
	"context"
	commonRepoPg "github.com/Azzonya/gophermart/internal/domain/common/repo/pg"
	"github.com/Azzonya/gophermart/internal/domain/order/model"
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

func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.Order, error) {
	var result []*model.Order
	var values []interface{}
	query := "SELECT * FROM order WHERE true"

	if pars.UserID != nil {
		query += " AND user_id = $1"
		values = append(values, *pars.UserID)
	}

	if pars.OrderNumber != nil {
		query += " AND code = $2"
		values = append(values, *pars.OrderNumber)
	}

	if pars.UploadedBefore != nil {
		query += " AND uploaded_at <= $3"
		values = append(values, *pars.UploadedBefore)
	}

	if pars.UploadedAfter != nil {
		query += " AND uploaded_at >= $4"
		values = append(values, *pars.UploadedAfter)
	}

	if pars.Status != nil {
		query += " AND status = $5"
		values = append(values, *pars.Status)
	}

	if pars.MinAccrual != nil {
		query += " AND accrual >= $6"
		values = append(values, *pars.MinAccrual)
	}

	if pars.MaxAccrual != nil {
		query += " AND accrual <= $7"
		values = append(values, *pars.MaxAccrual)
	}

	rows, err := r.Con.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var order *model.Order
		err = rows.Scan(&order.OrderNumber, &order.UploadedAt, &order.Status, &order.UserID, &order.Accrual)
		if err != nil {
			return nil, err
		}

		result = append(result, order)
	}

	return result, err
}

func (r *Repo) Create(ctx context.Context, obj *model.GetPars) error {
	_, err := r.Con.Exec(ctx, "INSERT INTO order (code, status, user_id, accrual) VALUES ($1, $2, $3, $4);", obj.OrderNumber, obj.Status, obj.UserID, obj.Accrual)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error) {
	var values []interface{}
	var result model.Order
	query := "SELECT * FROM order WHERE true"

	if len(pars.UserID) != 0 {
		query += " AND user_id = $1"
		values = append(values, pars.UserID)
	}

	if len(pars.OrderNumber) != 0 {
		query += " AND code = $2"
		values = append(values, pars.OrderNumber)
	}

	if len(pars.Status) != 0 {
		query += " AND status = $3"
		values = append(values, pars.Status)
	}

	if len(pars.UserID) != 0 {
		query += " AND user_id >= $4"
		values = append(values, pars.UserID)
	}

	if pars.Accrual != 0 {
		query += " AND accrual = $5"
		values = append(values, pars.Accrual)
	}

	err := r.Con.QueryRow(ctx, query, values...).Scan(&result)
	if err != nil {
		return nil, false, err
	}

	return &result, result.OrderNumber != "", err
}

func (r *Repo) Update(ctx context.Context, pars *model.GetPars) error {
	var values []interface{}

	query := "UPDATE order SET"

	if len(pars.UserID) != 0 {
		query += " AND user_id = $1"
		values = append(values, pars.UserID)
	}

	if len(pars.Status) != 0 {
		query += " AND status = $2"
		values = append(values, pars.Status)
	}

	if pars.Accrual != 0 {
		query += " AND accrual = $3"
		values = append(values, pars.Accrual)
	}

	_, err := r.Con.Exec(ctx, query, values...)

	return err
}

func (r *Repo) Delete(ctx context.Context, pars *model.GetPars) error {
	_, err := r.Con.Exec(ctx, "DELETE FROM order WHERE code = $1;", pars.OrderNumber)
	return err
}

func (r *Repo) Exists(ctx context.Context, orderNumber string) (bool, error) {
	var exist bool
	err := r.Con.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM order WHERE code = $1);", orderNumber).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, err
}
