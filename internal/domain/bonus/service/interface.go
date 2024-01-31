package service

import http_gw "github.com/Azzonya/gophermart/internal/client/accrual/http"

type RepoAccrualI interface {
	Send(orderNumber string) (*http_gw.RequestResult, error)
}
