package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azzonya/gophermart/internal/domain/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	Con *pgxpool.Pool
}

func New(con *pgxpool.Pool) *Repo {
	return &Repo{
		con,
	}
}

func (r *Repo) List(ctx context.Context, pars *user.ListPars) ([]*user.User, error) {
	var result []*user.User
	var values = make([]interface{}, 0)

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
		var userFound user.User
		err = rows.Scan(&userFound.ID, &userFound.Login, &userFound.Password, &userFound.Balance)
		if err != nil {
			return nil, err
		}

		result = append(result, &userFound)
	}

	return result, err
}

func (r *Repo) Create(ctx context.Context, obj *user.GetPars) error {
	_, err := r.Con.Exec(ctx, "INSERT INTO users (login, password) VALUES ($1, $2);", obj.Login, obj.Password)
	return err
}

func (r *Repo) Get(ctx context.Context, pars *user.GetPars) (*user.User, bool, error) {
	var values []interface{}
	var result user.User
	values = make([]interface{}, 0)

	query := "SELECT * FROM users WHERE true"

	paramNum := 1
	if len(pars.ID) != 0 {
		query += fmt.Sprintf(" AND id = $%d", paramNum)
		values = append(values, pars.ID)
		paramNum += 1
	}

	if len(pars.Login) != 0 {
		query += fmt.Sprintf(" AND login = $%d", paramNum)
		values = append(values, pars.Login)
		paramNum += 1
	}

	if pars.Balance != 0 {
		query += fmt.Sprintf(" AND balance = $%d", paramNum)
		values = append(values, pars.Balance)
		paramNum += 1
	}

	if len(pars.Password) != 0 {
		query += fmt.Sprintf(" AND password = $%d", paramNum)
		values = append(values, pars.Password)
		paramNum += 1
	}

	err := r.Con.QueryRow(ctx, query, values...).Scan(&result.ID, &result.Login, &result.Password, &result.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &result, result.Login != "", err
}

func (r *Repo) Update(ctx context.Context, pars *user.GetPars) error {
	var values []interface{}
	values = make([]interface{}, 0)

	query := "UPDATE users"

	paramNum := 1
	if len(pars.Login) != 0 {
		query += fmt.Sprintf(" SET login = $%d", paramNum)
		values = append(values, pars.Login)
		paramNum++
	}

	if len(pars.Password) != 0 {
		if len(values) > 0 {
			query += ","
		} else {
			query += " SET"
		}
		query += fmt.Sprintf(" password = $%d", paramNum)
		values = append(values, pars.Password)
		paramNum++
	}

	if pars.Balance != 0 {
		if len(values) > 0 {
			query += ","
		} else {
			query += " SET"
		}
		query += fmt.Sprintf(" balance = $%d", paramNum)
		values = append(values, pars.Balance)
		paramNum++
	}

	query += fmt.Sprintf(" WHERE id = '%s'", pars.ID)

	_, err := r.Con.Exec(ctx, query, values...)

	return err
}

func (r *Repo) Delete(ctx context.Context, pars *user.GetPars) error {
	_, err := r.Con.Exec(ctx, "DELETE FROM users WHERE login = $1;", pars.Login)
	return err
}

func (r *Repo) Exists(ctx context.Context, login string) (bool, error) {
	var exist bool

	err := r.Con.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM users WHERE login = $1);", login).Scan(&exist)
	if err != nil {
		return false, err
	}

	return exist, err
}
