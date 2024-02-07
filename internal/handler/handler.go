package handler

import (
	"github.com/Azzonya/gophermart/internal/domain/auth"
	"github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"github.com/Azzonya/gophermart/internal/usecase/user"
)

type UserHandlers struct {
	auth                     *auth.Auth
	userUsecase              *user.Usecase
	orderUsecase             *order.Usecase
	bonusTransactionsUsecase *bonustransactions.Usecase

	pgDsn string
}

func New(auth *auth.Auth, userUsecase *user.Usecase, orderUsecase *order.Usecase, bonusTransactionsUsecase *bonustransactions.Usecase, pgDsn string) *UserHandlers {
	return &UserHandlers{
		auth:                     auth,
		userUsecase:              userUsecase,
		orderUsecase:             orderUsecase,
		bonusTransactionsUsecase: bonusTransactionsUsecase,
		pgDsn:                    pgDsn,
	}
}
