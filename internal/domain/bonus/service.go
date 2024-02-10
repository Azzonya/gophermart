package bonus

import (
	"context"
	"fmt"
	http_gw "github.com/Azzonya/gophermart/internal/client/accrual/http"
	bonusTransactionsModel "github.com/Azzonya/gophermart/internal/domain/bonustransactions"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order"
	userModel "github.com/Azzonya/gophermart/internal/domain/user"
	"github.com/Azzonya/gophermart/internal/storage"
	"github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"github.com/Azzonya/gophermart/internal/usecase/user"
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
	userService              user.UserServiceI
	wg                       sync.WaitGroup
}

func New(accrual ClientAccrualI, bonusTransactionsService bonustransactions.WithdrawalServiceI, orderService order.OrderServiceI, user user.UserServiceI) *Service {
	return &Service{
		accrual:                  accrual,
		bonusTransactionsService: bonusTransactionsService,
		orderService:             orderService,
		userService:              user,
	}
}

func (s *Service) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

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

		foundUser, _, err := s.userService.Get(ctx, &userModel.GetPars{
			ID: v.UserID,
		})
		if err != nil {
			slog.Error("failed to get user")
			return
		}

		err = s.userService.Update(ctx, &userModel.GetPars{
			ID:      v.UserID,
			Balance: foundUser.Balance + responseResult.Accrual,
		})
		if err != nil {
			slog.Error("failed to update user balance")
			return
		}
	}
}

func (s *Service) WithdrawBalance(ctx context.Context, pars *bonusTransactionsModel.GetPars) error {
	foundUser, _, err := s.userService.Get(ctx, &userModel.GetPars{
		ID: pars.UserID,
	})
	if err != nil {
		return err
	}

	if foundUser.Balance < pars.Sum {
		return storage.ErrUserInsufficientBalance{}
	}

	_, orderExist, err := s.orderService.Get(ctx, &orderModel.GetPars{
		OrderNumber: pars.OrderNumber,
	})
	if err != nil {
		return err
	}
	if !orderExist {
		return storage.ErrOrderNotExist{OrderNumber: pars.OrderNumber}
	}

	err = s.bonusTransactionsService.Create(ctx, pars)
	if err != nil {
		return err
	}

	err = s.userService.Update(ctx, &userModel.GetPars{
		ID:      foundUser.ID,
		Balance: foundUser.Balance - pars.Sum,
	})
	if err != nil {
		return err
	}

	return err
}
