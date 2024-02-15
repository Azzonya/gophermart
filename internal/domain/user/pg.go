package user

import (
	"context"
	"errors"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/Azzonya/gophermart/internal/errs"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repoDBI interface {
	ListUsers(ctx context.Context, pars *entities.UserListPars) ([]*entities.User, error)
	Create(ctx context.Context, obj *entities.User) error
	Get(ctx context.Context, pars *entities.UserParameters) (*entities.User, error)
	Update(ctx context.Context, pars *entities.UserParameters) error
	Delete(ctx context.Context, pars *entities.UserParameters) error
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

func (r *Repo) ListUsers(ctx context.Context, pars *entities.UserListPars) ([]*entities.User, error) {
	queryBuilder := squirrel.Select("*").From("users").Where(squirrel.Eq{"true": true})

	if pars.Login != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"login": *pars.Login})
	}

	if pars.MaxBalance != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"balance": *pars.MaxBalance})
	}

	if pars.MinBalance != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"balance": *pars.MinBalance})
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

	var result []*entities.User

	for rows.Next() {
		var userFound entities.User
		err = rows.Scan(&userFound.ID, &userFound.Login, &userFound.Password, &userFound.Balance)
		if err != nil {
			return nil, err
		}
		result = append(result, &userFound)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repo) Create(ctx context.Context, obj *entities.User) error {
	queryBuilder := squirrel.Insert("users").
		Columns("login", "password").
		Values(obj.Login, obj.Password)

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == pgerrcode.UniqueViolation {
			return errs.ErrUserNotUniq{Login: obj.Login}
		}
		return err
	}

	return nil
}

func (r *Repo) Get(ctx context.Context, pars *entities.UserParameters) (*entities.User, error) {
	queryBuilder := squirrel.Select("*").From("users").Where(squirrel.Expr("true"))

	if len(pars.ID) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})
	}
	if len(pars.Login) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"login": pars.Login})
	}
	if pars.Balance != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"balance": pars.Balance})
	}
	if len(pars.Password) != 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"password": pars.Password})
	}

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var result entities.User
	err = r.Con.QueryRow(ctx, sql, args...).Scan(&result.ID, &result.Login, &result.Password, &result.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &result, nil
}

func (r *Repo) Update(ctx context.Context, pars *entities.UserParameters) error {
	queryBuilder := squirrel.Update("users")

	if len(pars.Login) != 0 {
		queryBuilder = queryBuilder.Set("login", pars.Login)
	}
	if len(pars.Password) != 0 {
		queryBuilder = queryBuilder.Set("password", pars.Password)
	}
	if pars.Balance != 0 {
		queryBuilder = queryBuilder.Set("balance", pars.Balance)
	}

	queryBuilder = queryBuilder.Where(squirrel.Eq{"id": pars.ID})

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

func (r *Repo) Delete(ctx context.Context, pars *entities.UserParameters) error {
	queryBuilder := squirrel.Delete("users")

	queryBuilder = queryBuilder.Where(squirrel.Eq{"login": pars.Login})

	sql, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	_, err = r.Con.Exec(ctx, sql, args...)
	return err
}

func (r *Repo) Exists(ctx context.Context, login string) (bool, error) {
	existsQuery := squirrel.Select("SELECT EXISTS (SELECT 1 FROM users WHERE login = ?);", login)

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
