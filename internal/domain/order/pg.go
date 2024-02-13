package order

import (
	"context"
	"errors"
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

func (r *Repo) List(ctx context.Context, pars *ListPars) ([]*Order, error) {
	queryBuilder := squirrel.Select("code", "uploaded_at", "status", "user_id").From("orders").Where(squirrel.Eq{"true": true})

	var values = make(map[string]interface{})

	if pars.UserID != nil {
		values["user_id"] = *pars.UserID
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": values["user_id"]})
	}

	if pars.OrderNumber != nil {
		queryBuilder = queryBuilder.Where(squirrel.Or{squirrel.Eq{"code": *pars.OrderNumber}, squirrel.Eq{"$null": nil}})
	}

	if pars.UploadedBefore != nil {
		values["uploaded_at"] = *pars.UploadedBefore
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"uploaded_at": values["uploaded_at"]})
	}

	if pars.UploadedAfter != nil {
		values["uploaded_at"] = *pars.UploadedAfter
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"uploaded_at": values["uploaded_at"]})
	}

	if pars.Status != nil {
		values["status"] = *pars.Status
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": values["status"]})
	}

	if pars.Statuses != nil && len(pars.Statuses) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Statuses})
	}

	if len(pars.OrderBy) != 0 {
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("uploaded_at %s", pars.OrderBy))
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.Con.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*Order
	for rows.Next() {
		var orderFound Order
		err = rows.Scan(&orderFound.OrderNumber, &orderFound.UploadedAt, &orderFound.Status, &orderFound.UserID)
		if err != nil {
			return nil, err
		}

		result = append(result, &orderFound)
	}

	return result, nil
}

func (r *Repo) Create(ctx context.Context, obj *Order) error {
	insert := squirrel.Insert("orders").
		Columns("code", "status", "user_id").
		Values(obj.OrderNumber, obj.Status, obj.UserID)

	query, args, err := insert.ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, query, args...)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *GetPars) (*Order, error) {
	var result Order

	queryBuilder := squirrel.Select("*").From("orders").Where(squirrel.Eq{"true": true})

	if len(pars.UserID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	if len(pars.OrderNumber) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"code": pars.OrderNumber})
	}

	if len(pars.Status) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Status})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.OrderNumber, &result.UploadedAt, &result.Status, &result.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *Repo) Update(ctx context.Context, pars *GetPars) error {
	queryBuilder := squirrel.Update("orders")

	if len(pars.UserID) != 0 {
		queryBuilder = queryBuilder.Set("user_id", pars.UserID)
	}

	if len(pars.Status) != 0 {
		queryBuilder = queryBuilder.Set("status", pars.Status)
	}

	queryBuilder = queryBuilder.Where(squirrel.Eq{"code": pars.OrderNumber})

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

func (r *Repo) Delete(ctx context.Context, pars *GetPars) error {
	queryBuilder := squirrel.Delete("orders")

	queryBuilder = queryBuilder.Where(squirrel.Eq{"code": pars.OrderNumber})

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

func (r *Repo) Exists(ctx context.Context, orderNumber string) (bool, error) {
	existsQuery := squirrel.Select("SELECT EXISTS (SELECT 1 FROM orders WHERE code = $1);", orderNumber)

	query, args, err := existsQuery.ToSql()
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
