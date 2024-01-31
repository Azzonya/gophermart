package service

import http_gw "github.com/Azzonya/gophermart/internal/client/accrual/http-gw"

type RepoAccrualI interface {
	Send(orderNumber string) (*http_gw.RequestResult, error)
}
