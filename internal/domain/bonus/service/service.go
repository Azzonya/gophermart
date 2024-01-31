package service

import (
	"context"
	"fmt"
	bonus_transactionsModel "github.com/Azzonya/gophermart/internal/domain/bonus_transactions/model"
	orderModel "github.com/Azzonya/gophermart/internal/domain/order/model"
	"github.com/Azzonya/gophermart/internal/usecase/bonus_transactions"
	"github.com/Azzonya/gophermart/internal/usecase/order"
	"sync"
	"time"
)

type Service struct {
	repoAccrual              RepoAccrualI
	bonusTransactionsService bonus_transactions.WithdrawalServiceI
	orderService             order.OrderServiceI
	Wg                       sync.WaitGroup
}

func New(repoAccrual RepoAccrualI, bonusTransactionsService bonus_transactions.WithdrawalServiceI, orderService order.OrderServiceI) *Service {
	return &Service{
		repoAccrual:              repoAccrual,
		bonusTransactionsService: bonusTransactionsService,
		orderService:             orderService,
	}
}

func (s *Service) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)

	// Функция, которая будет вызываться при каждом тике таймера
	go func() {
		defer ticker.Stop() // Отложенная остановка таймера при завершении горутины
		for {
			select {
			case <-ticker.C:
				// Увеличиваем счетчик горутин перед созданием новой
				s.Wg.Add(1)
				// Вызываем функцию отправки GET-запроса
				go s.UpdateAccrualInfo(context.Background())
			}
		}
	}()

	defer ticker.Stop()
}

func (s *Service) UpdateAccrualInfo(ctx context.Context) error {
	fmt.Println(123)
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
		responseResult, err := s.repoAccrual.Send(v.OrderNumber)
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
		}

		bonus_transaction, exist, err := s.bonusTransactionsService.Get(ctx, &bonus_transactionsModel.GetPars{
			OrderNumber:     v.OrderNumber,
			TransactionType: bonus_transactionsModel.Accrual,
		})
		if exist {
			if bonus_transaction.Sum != responseResult.Accrual {
				err = s.bonusTransactionsService.Update(ctx, &bonus_transactionsModel.GetPars{
					OrderNumber: v.OrderNumber,
					Sum:         responseResult.Accrual,
				})
				if err != nil {
					return err
				}
			}
			continue
		}

		err = s.bonusTransactionsService.Create(ctx, &bonus_transactionsModel.GetPars{
			OrderNumber:     v.OrderNumber,
			UserID:          v.UserID,
			TransactionType: bonus_transactionsModel.Accrual,
			Sum:             responseResult.Accrual,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
