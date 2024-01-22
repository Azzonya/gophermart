package db

import (
	"context"
	commonRepoPg "github.com/Azzonya/gophermart/internal/domain/common/repo/pg"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
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

func (r *Repo) List(ctx context.Context, pars *model.ListPars) ([]*model.User, error) {
	var result []*model.User
	var values []interface{}
	query := "SELECT * FROM users WHERE true"

	if pars.Login != nil {
		query += " AND login = $1"
		values = append(values, *pars.Login)
	}

	if pars.MaxBalance != nil {
		query += " AND balance <= $2"
		values = append(values, *pars.MaxBalance)
	}

	if pars.MinBalance != nil {
		query += " AND balance <= $3"
		values = append(values, *pars.MinBalance)
	}

	rows, err := r.Con.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var user *model.User
		err = rows.Scan(&user.ID, &user.Login, &user.Password, &user.Balance)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, err
}

func (r *Repo) Create(ctx context.Context, obj *model.GetPars) error {
	_, err := r.Con.Exec(ctx, "INSERT INTO user (login, password) VALUES ($1, $2);", obj.Login, obj.Password)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	var values []interface{}
	var result model.User
	query := "SELECT * FROM user WHERE true"

	if len(pars.Login) != 0 {
		query += " AND login = $1"
		values = append(values, pars.Login)
	}

	if pars.Balance != 0 {
		query += " AND balance = $2"
		values = append(values, pars.Balance)
	}

	if len(pars.Password) != 0 {
		query += " AND password = $3"
		values = append(values, pars.Password)
	}

	err := r.Con.QueryRow(ctx, query, values...).Scan(&result)
	if err != nil {
		return nil, false, err
	}

	return &result, result.Login != "", err
}

func (r *Repo) Update(ctx context.Context, pars *model.GetPars) error {
	var values []interface{}

	query := "UPDATE user SET"

	if len(pars.Login) != 0 {
		query += " AND login = $1"
		values = append(values, pars.Login)
	}

	if len(pars.Password) != 0 {
		query += " AND password = $2"
		values = append(values, pars.Password)
	}

	if pars.Balance != 0 {
		query += " AND balance = $3"
		values = append(values, pars.Balance)
	}

	_, err := r.Con.Exec(ctx, query, values...)

	return err
}

func (r *Repo) Delete(ctx context.Context, pars *model.GetPars) error {
	_, err := r.Con.Exec(ctx, "DELETE FROM user WHERE login = $1;", pars.Login)
	return err
}

func (r *Repo) Exists(ctx context.Context, login string) (bool, error) {
	var exist bool
	err := r.Con.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM user WHERE login = $1);", login).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, err
}
