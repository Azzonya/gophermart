package accrual

import (
	http_gw "github.com/Azzonya/gophermart/internal/client/accrual/http"
)

type Client interface {
	Send(orderNumber string) (*http_gw.RequestResult, error)
}
