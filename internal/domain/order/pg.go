package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepoDBI interface {
	ListOrders(ctx context.Context, pars *entities.OrderListPars) ([]*entities.Order, error)
	Create(ctx context.Context, obj *entities.Order) error
	Get(ctx context.Context, pars *entities.OrderParameters) (*entities.Order, error)
	Update(ctx context.Context, pars *entities.OrderParameters) error
	Delete(ctx context.Context, pars *entities.OrderParameters) error
	Exists(ctx context.Context, orderNumber string) (bool, error)
}

type Repo struct {
	Con *pgxpool.Pool
}

func NewRepo(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

func (r *Repo) ListOrders(ctx context.Context, pars *entities.OrderListPars) ([]*entities.Order, error) {
	queryBuilder := squirrel.Select("code", "uploaded_at", "status", "user_id").From("orders").Where(squirrel.Eq{"true": true})

	if pars.UserID != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	if pars.OrderNumber != nil {
		queryBuilder = queryBuilder.Where(squirrel.Or{squirrel.Eq{"code": *pars.OrderNumber}, squirrel.Eq{"$null": nil}})
	}

	if pars.UploadedBefore != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"uploaded_at": pars.UploadedBefore})
	}

	if pars.UploadedAfter != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"uploaded_at": pars.UploadedAfter})
	}

	if pars.Status != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Status})
	}

	if pars.Statuses != nil && len(pars.Statuses) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Statuses})
	}

	if len(pars.OrderBy) != 0 {
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("uploaded_at %s", pars.OrderBy))
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

	var result []*entities.Order
	for rows.Next() {
		var orderFound entities.Order
		err = rows.Scan(&orderFound.OrderNumber, &orderFound.UploadedAt, &orderFound.Status, &orderFound.UserID)
		if err != nil {
			return nil, err
		}

		result = append(result, &orderFound)
	}

	return result, nil
}

func (r *Repo) Create(ctx context.Context, obj *entities.Order) error {
	insert := squirrel.Insert("orders").
		Columns("code", "status", "user_id").
		Values(obj.OrderNumber, obj.Status, obj.UserID)

	query, args, err := insert.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, query, args...)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *entities.OrderParameters) (*entities.Order, error) {
	var result entities.Order

	queryBuilder := squirrel.Select("*").From("orders")

	if len(pars.UserID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": pars.UserID})
	}

	if len(pars.OrderNumber) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"code": pars.OrderNumber})
	}

	if len(pars.Status) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"status": pars.Status})
	}

	queryBuilder = queryBuilder.Limit(1)

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
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

func (r *Repo) Update(ctx context.Context, pars *entities.OrderParameters) error {
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

func (r *Repo) Delete(ctx context.Context, pars *entities.OrderParameters) error {
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
