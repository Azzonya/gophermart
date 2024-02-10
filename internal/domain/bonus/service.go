package bonus

import (
	"context"
	"fmt"
	http_gw "github.com/Azzonya/gophermart/internal/client/accrual/http"
	bonusTransactionsModel "github.com/Azzonya/gophermart/internal/domain/bonustransactions"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order"
	"github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"log/slog"
	"sync"
	"time"
)

type ClientAccrualI interface {
	Send(orderNumber string) (*http_gw.RequestResult, error)
}

type Service struct {
	accrual                  ClientAccrualI
	bonusTransactionsService bonustransactions.WithdrawalServiceI
	orderService             order.OrderServiceI
	wg                       sync.WaitGroup
}

func New(accrual ClientAccrualI, bonusTransactionsService bonustransactions.WithdrawalServiceI, orderService order.OrderServiceI) *Service {
	return &Service{
		accrual:                  accrual,
		bonusTransactionsService: bonusTransactionsService,
		orderService:             orderService,
	}
}

func (s *Service) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.wg.Add(1)
				s.updateAccrualInfo(ctx)
			}
		}
	}()
}

func (s *Service) Wait() {
	s.wg.Wait()
}

func (s *Service) updateAccrualInfo(ctx context.Context) {
	defer s.wg.Done()

	statuses := []orderModel.OrderStatus{
		orderModel.OrderStatusNew,
		orderModel.OrderStatusProcessing,
	}
	orders, err := s.orderService.List(ctx, &orderModel.ListPars{
		Statuses: statuses,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("failed to fetch orders: %s", err.Error()))
		return
	}

	for _, v := range orders {

		responseResult, err := s.accrual.Send(v.OrderNumber)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to send accrual for order %s: %s", v.OrderNumber, err.Error()))
			return
		}

		if responseResult == nil {
			return
		}

		if string(v.Status) != responseResult.Status {
			err = s.orderService.Update(ctx, &orderModel.GetPars{
				Status:      orderModel.OrderStatus(responseResult.Status),
				OrderNumber: v.OrderNumber,
			})
			if err != nil {
				slog.Error(fmt.Sprintf("failed to update order status for order %s: %s", v.OrderNumber, err.Error()))
				return
			}
		}

		_, exist, err := s.bonusTransactionsService.Get(ctx, &bonusTransactionsModel.GetPars{
			OrderNumber:     v.OrderNumber,
			TransactionType: bonusTransactionsModel.Accrual,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("failed to check bonus transaction existence for order %s: %s", v.OrderNumber, err.Error()))
			return
		}
		if exist {
			slog.Error(fmt.Sprintf("bonus transaction with order number %s already exists", v.OrderNumber))
			return
		}

		err = s.bonusTransactionsService.Create(ctx, &bonusTransactionsModel.GetPars{
			OrderNumber:     v.OrderNumber,
			UserID:          v.UserID,
			TransactionType: bonusTransactionsModel.Accrual,
			Sum:             responseResult.Accrual,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("failed to create bonus transaction for order %s: %s", v.OrderNumber, err.Error()))
			return
		}
	}
}
