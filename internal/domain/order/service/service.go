package service

import (
	"context"
	bonus_transactionsModel "github.com/Azzonya/gophermart/internal/domain/bonusTransactions/model"
	"github.com/Azzonya/gophermart/internal/domain/order/model"
	bonus_transactions "github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"strconv"
	"time"
)

type Service struct {
	repoDB                   RepoDBI
	bonusTransactionsService bonus_transactions.WithdrawalServiceI
}

func New(repoDB RepoDBI, bonusTransactionsService bonus_transactions.WithdrawalServiceI) *Service {
	return &Service{
		repoDB:                   repoDB,
		bonusTransactionsService: bonusTransactionsService,
	}
}

func (s *Service) IsLuhnValid(orderNumber string) bool {
	// Преобразование строки в массив цифр
	digits := make([]int, len(orderNumber))

	for i, char := range orderNumber {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Применение алгоритма Луна
	total := 0
	for i := len(digits) - 2; i >= 0; i -= 2 {
		digits[i] *= 2
		if digits[i] > 9 {
			digits[i] -= 9
		}
	}

	for _, digit := range digits {
		total += digit
	}

	return total%10 == 0
}

func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.Order, error) {
	return s.repoDB.List(ctx, pars)
}

func (s *Service) ListWithAccrual(ctx context.Context, pars *model.ListPars) ([]*model.OrderWithAccrual, error) {
	orders, err := s.repoDB.List(ctx, pars)
	if err != nil {
		return nil, err
	}

	orderMap := make(map[string]model.Order)
	for _, order := range orders {
		orderMap[order.OrderNumber] = *order
	}

	bonusTransactions, err := s.bonusTransactionsService.List(ctx, &bonus_transactionsModel.ListPars{
		UserID:          pars.UserID,
		TransactionType: bonus_transactionsModel.Accrual,
	})
	if err != nil {
		return nil, err
	}

	bonusMap := make(map[string]bonus_transactionsModel.BonusTransaction)
	for _, bonus := range bonusTransactions {
		bonusMap[bonus.OrderNumber] = *bonus
	}

	var result []*model.OrderWithAccrual
	for _, order := range orders {
		bonusTransaction, exists := bonusMap[order.OrderNumber]

		accrualSum := 0
		if exists {
			accrualSum = bonusTransaction.Sum
		}

		result = append(result, &model.OrderWithAccrual{
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			Accrual:     accrualSum,
			UploadedAt:  order.UploadedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *Service) Create(ctx context.Context, obj *model.GetPars) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error) {
	return s.repoDB.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *model.GetPars) error {
	return s.repoDB.Update(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDB.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDB.Exists(ctx, orderNumber)
}
