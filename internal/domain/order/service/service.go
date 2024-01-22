package service

import (
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
