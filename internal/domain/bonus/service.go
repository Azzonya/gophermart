package bonus

import (
	"context"
	"fmt"
	http_gw "github.com/Azzonya/gophermart/internal/client/accrual/http"
	bonusTransactionsModel "github.com/Azzonya/gophermart/internal/domain/bonusTransactions"
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

func (s *Service) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()
		for range time.Tick(interval) {
			s.wg.Add(1)

			go func() {
				err := s.updateAccrualInfo(context.Background())
				if err != nil {
					slog.Error(err.Error())
				}
			}()
		}
	}()

	defer ticker.Stop()
}

func (s *Service) Wait() {
	s.wg.Wait()
}

func (s *Service) updateAccrualInfo(ctx context.Context) error {
	statuses := []orderModel.OrderStatus{
		orderModel.OrderStatusNew,
		orderModel.OrderStatusProcessing,
	}
	orders, err := s.orderService.List(ctx, &orderModel.ListPars{
		Statuses: statuses,
	})
	if err != nil {
		return err
	}

	for _, v := range orders {
		responseResult, err := s.accrual.Send(v.OrderNumber)
		if err != nil {
			return err
		}

		if responseResult == nil {
			continue
		}

		if string(v.Status) != responseResult.Status {
			err = s.orderService.Update(ctx, &orderModel.GetPars{
				Status:      orderModel.OrderStatus(responseResult.Status),
				OrderNumber: v.OrderNumber,
			})
			if err != nil {
				return err
			}
		}

		_, exist, err := s.bonusTransactionsService.Get(ctx, &bonusTransactionsModel.GetPars{
			OrderNumber:     v.OrderNumber,
			TransactionType: bonusTransactionsModel.Accrual,
		})
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("bonus transaction with this order number exist")
		}

		err = s.bonusTransactionsService.Create(ctx, &bonusTransactionsModel.GetPars{
			OrderNumber:     v.OrderNumber,
			UserID:          v.UserID,
			TransactionType: bonusTransactionsModel.Accrual,
			Sum:             responseResult.Accrual,
		})
		if err != nil {
			return err
		}
	}

	s.wg.Done()
	return nil
}
