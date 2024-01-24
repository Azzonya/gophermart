package service

import (
	"context"
	"github.com/Azzonya/gophermart/internal/domain/order/model"
	"strconv"
)

type Service struct {
	repoDb RepoDbI
}

func (s *Service) isLuhnValid(orderNumber string) bool {
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
	return s.repoDb.List(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *model.GetPars) error {
	return s.repoDb.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.Order, bool, error) {
	return s.repoDb.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *model.GetPars, obj *model.GetPars) error {
	return s.repoDb.Update(ctx, pars, obj)
}

func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDb.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDb.Exists(ctx, orderNumber)
}
