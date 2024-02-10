package storage

import (
	"fmt"
)

type ErrUserNotUniq struct {
	Login string
}

func (err ErrUserNotUniq) Error() string {
	return fmt.Sprintf("user with login \"%s\" already exists", err.Login)
}

type ErrOrderNumberLuhnValid struct {
	OrderNumber string
}

func (err ErrOrderNumberLuhnValid) Error() string {
	return fmt.Sprintf("order with number - \"%s\" not Luhn valie", err.OrderNumber)
}

type ErrOrderUploaded struct {
	OrderNumber string
}

func (err ErrOrderUploaded) Error() string {
	return fmt.Sprintf("order with number - \"%s\" uploaded", err.OrderNumber)
}

type ErrOrderUploadedByAnotherUser struct {
	OrderNumber string
}

func (err ErrOrderUploadedByAnotherUser) Error() string {
	return fmt.Sprintf("order with number - \"%s\" uploaded by another user", err.OrderNumber)
}

type ErrUserInsufficientBalance struct {
}

func (err ErrUserInsufficientBalance) Error() string {
	return "Insufficient balance"
}
