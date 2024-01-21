package common

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Base struct {
	Con *pgxpool.Pool
}

func NewBase(con *pgxpool.Pool) *Base {
	return &Base{
		Con: con,
	}
}
