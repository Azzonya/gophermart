package order

import (
	"context"
	"github.com/Azzonya/gophermart/internal/entities"
	bonusTransactions "github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"strconv"
	"time"
)

type Service struct {
	repoDB                   OrderRepoDBI
	bonusTransactionsService bonusTransactions.WithdrawalServiceI
}

func New(repoDB OrderRepoDBI, bonusTransactionsService bonusTransactions.WithdrawalServiceI) *Service {
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

func (s *Service) List(ctx context.Context, pars *entities.OrderListPars) ([]*entities.Order, error) {
	return s.repoDB.ListOrders(ctx, pars)
}

func (s *Service) ListWithAccrual(ctx context.Context, pars *entities.OrderListPars) ([]*entities.OrderWithAccrual, error) {
	orders, err := s.repoDB.ListOrders(ctx, pars)
	if err != nil {
		return nil, err
	}

	orderMap := make(map[string]entities.Order)
	for _, order := range orders {
		orderMap[order.OrderNumber] = *order
	}

	bonusTransactionsList, err := s.bonusTransactionsService.List(ctx, &entities.BonusTransactionsListPars{
		UserID:          pars.UserID,
		TransactionType: entities.Accrual,
	})
	if err != nil {
		return nil, err
	}

	bonusMap := make(map[string]entities.BonusTransaction)
	for _, bonus := range bonusTransactionsList {
		bonusMap[bonus.OrderNumber] = *bonus
	}

	var result []*entities.OrderWithAccrual
	var bonusTransaction entities.BonusTransaction
	var accrualSum float32
	var exists bool

	for _, order := range orders {
		bonusTransaction, exists = bonusMap[order.OrderNumber]

		accrualSum = 0
		if exists {
			accrualSum = bonusTransaction.Sum
		}

		result = append(result, &entities.OrderWithAccrual{
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			Accrual:     accrualSum,
			UploadedAt:  order.UploadedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

func (s *Service) Create(ctx context.Context, obj *entities.Order) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *entities.OrderParameters) (*entities.Order, error) {
	return s.repoDB.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *entities.OrderParameters) error {
	return s.repoDB.Update(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, pars *entities.OrderParameters) error {
	return s.repoDB.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDB.Exists(ctx, orderNumber)
}
