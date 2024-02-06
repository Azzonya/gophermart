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
}

func New(auth *auth.Auth, userUsecase *user.Usecase, orderUsecase *order.Usecase, bonusTransactionsUsecase *bonustransactions.Usecase) *UserHandlers {
	return &UserHandlers{
		auth:                     auth,
		userUsecase:              userUsecase,
		orderUsecase:             orderUsecase,
		bonusTransactionsUsecase: bonusTransactionsUsecase,
	}
}
