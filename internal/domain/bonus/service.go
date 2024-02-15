package bonus

import (
	"context"
	"fmt"
	http_gw "github.com/Azzonya/gophermart/internal/client/accrual/http"
	"github.com/Azzonya/gophermart/internal/entities"
	"github.com/Azzonya/gophermart/internal/errs"
	"log/slog"
	"sync"
	"time"
)

type WithdrawalServiceI interface {
	Get(ctx context.Context, pars *entities.BonusTransactionsParameters) (*entities.BonusTransaction, error)
	Create(ctx context.Context, obj *entities.BonusTransaction) error
}

type OrderServiceI interface {
	List(ctx context.Context, pars *entities.OrderListPars) ([]*entities.Order, error)
	Update(ctx context.Context, pars *entities.OrderParameters) error
}

type UserServiceI interface {
	Get(ctx context.Context, pars *entities.UserParameters) (*entities.User, error)
	Update(ctx context.Context, pars *entities.UserParameters) error
}

type ClientAccrualI interface {
	Send(orderNumber string) (*http_gw.RequestResult, error)
}

type Service struct {
	accrual                  ClientAccrualI
	bonusTransactionsService WithdrawalServiceI
	orderService             OrderServiceI
	userService              UserServiceI
	wg                       sync.WaitGroup
}

func New(accrual ClientAccrualI, bonusTransactionsService WithdrawalServiceI, orderService OrderServiceI, user UserServiceI) *Service {
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

	statuses := []entities.OrderStatus{
		entities.OrderStatusNew,
		entities.OrderStatusProcessing,
	}
	orders, err := s.orderService.List(ctx, &entities.OrderListPars{
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
			err = s.orderService.Update(ctx, &entities.OrderParameters{
				Status:      entities.OrderStatus(responseResult.Status),
				OrderNumber: v.OrderNumber,
			})
			if err != nil {
				slog.Error(fmt.Sprintf("failed to update order status for order %s: %s", v.OrderNumber, err.Error()))
				return
			}
		}

		bonusTransaction, err := s.bonusTransactionsService.Get(ctx, &entities.BonusTransactionsParameters{
			OrderNumber:     v.OrderNumber,
			TransactionType: entities.Accrual,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("failed to check bonus transaction existence for order %s: %s", v.OrderNumber, err.Error()))
			return
		}
		if bonusTransaction != nil {
			slog.Error(fmt.Sprintf("bonus transaction with order number %s already exists", v.OrderNumber))
			return
		}

		err = s.bonusTransactionsService.Create(ctx, &entities.BonusTransaction{
			OrderNumber:     v.OrderNumber,
			UserID:          v.UserID,
			TransactionType: entities.Accrual,
			Sum:             responseResult.Accrual,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("failed to create bonus transaction for order %s: %s", v.OrderNumber, err.Error()))
			return
		}

		foundUser, err := s.userService.Get(ctx, &entities.UserParameters{
			ID: v.UserID,
		})
		if err != nil {
			slog.Error("failed to get user")
			return
		}

		err = s.userService.Update(ctx, &entities.UserParameters{
			ID:      v.UserID,
			Balance: foundUser.Balance + responseResult.Accrual,
		})
		if err != nil {
			slog.Error("failed to update user balance")
			return
		}
	}
}

func (s *Service) WithdrawBalance(ctx context.Context, pars *entities.BonusTransaction) error {
	foundUser, err := s.userService.Get(ctx, &entities.UserParameters{
		ID: pars.UserID,
	})
	if err != nil {
		return err
	}

	if foundUser.Balance < pars.Sum {
		return errs.ErrUserInsufficientBalance{}
	}

	err = s.bonusTransactionsService.Create(ctx, pars)
	if err != nil {
		return err
	}

	err = s.userService.Update(ctx, &entities.UserParameters{
		ID:      foundUser.ID,
		Balance: foundUser.Balance - pars.Sum,
	})
	if err != nil {
		return err
	}

	return err
}
