package service

import (
	"context"
	bonus_transactionsModel "github.com/Azzonya/gophermart/internal/domain/bonusTransactions/model"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
	"github.com/Azzonya/gophermart/internal/usecase/bonustransactions"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repoDB                   repoDBI
	bonusTransactionsService bonustransactions.WithdrawalServiceI
}

func New(repoDB repoDBI, bonusTransactionsService bonustransactions.WithdrawalServiceI) *Service {
	return &Service{
		repoDB:                   repoDB,
		bonusTransactionsService: bonusTransactionsService,
	}
}

func (s *Service) IsValidPassword(password string, plainPassword string) bool {
	// Сравниваем хэшированный пароль из базы данных с переданным паролем
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(plainPassword))
	return err == nil
}

func (s *Service) HashPassword(password string) (string, error) {
	// Хэшируем пароль перед сохранением в базу данных
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *Service) IsLoginTaken(ctx context.Context, login string) (bool, error) {
	return s.repoDB.Exists(ctx, login)
}

func (s *Service) Register(ctx context.Context, user *model.GetPars) (*model.User, error) {
	var err error
	user.Password, err = s.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	err = s.repoDB.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	u, _, err := s.repoDB.Get(ctx, &model.GetPars{
		Login: user.Login,
	})
	if err != nil {
		return nil, err
	}

	return u, err
}

func (s *Service) GetBalanceWithWithdrawn(ctx context.Context, pars *model.GetPars) (*model.UserBalance, error) {
	user, _, err := s.Get(ctx, pars)
	if err != nil {
		return nil, err
	}

	bonusTransactionsList, err := s.bonusTransactionsService.List(ctx, &bonus_transactionsModel.ListPars{
		UserID:          &pars.ID,
		TransactionType: bonus_transactionsModel.Debit,
	})
	if err != nil {
		return nil, err
	}

	withdrawn := 0
	for _, v := range bonusTransactionsList {
		withdrawn += v.Sum
	}

	return &model.UserBalance{
		Current:   user.Balance,
		Withdrawn: withdrawn,
	}, nil
}

func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.User, error) {
	return s.repoDB.List(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *model.GetPars) error {
	return s.repoDB.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	return s.repoDB.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *model.GetPars) error {
	return s.repoDB.Update(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDB.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDB.Exists(ctx, orderNumber)
}
