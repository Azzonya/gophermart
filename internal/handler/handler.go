package handler

import (
	"github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"github.com/Azzonya/gophermart/internal/usecase/user"
)

type UserHandlers struct {
	userUsecase              *user.Usecase
	orderUsecase             *order.Usecase
	bonusTransactionsUsecase *bonustransactions.Usecase

	pgDsn string
}

func New(userUsecase *user.Usecase, orderUsecase *order.Usecase, bonusTransactionsUsecase *bonustransactions.Usecase, pgDsn string) *UserHandlers {
	return &UserHandlers{
		userUsecase:              userUsecase,
		orderUsecase:             orderUsecase,
		bonusTransactionsUsecase: bonusTransactionsUsecase,
		pgDsn:                    pgDsn,
	}
}
