package accrual

import "context"

type Client interface {
	Send(ctx context.Context, obj *SendReq) ([]byte, error)
}
