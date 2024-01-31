package accrual

import (
	"github.com/Azzonya/gophermart/internal/client/accrual"
)

type Repo struct {
	client accrual.Client
}

func New(client accrual.Client) *Repo {
	return &Repo{
		client: client,
	}
}
