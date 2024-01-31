package service

import (
	"context"
	bonus_transactionsModel "github.com/Azzonya/gophermart/internal/domain/bonus_transactions/model"
	"github.com/Azzonya/gophermart/internal/domain/user/model"
	"github.com/Azzonya/gophermart/internal/usecase/bonus_transactions"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repoDb                   RepoDbI
	bonusTransactionsService bonus_transactions.WithdrawalServiceI
}

func New(repoDb RepoDbI, bonusTransactionsService bonus_transactions.WithdrawalServiceI) *Service {
	return &Service{
		repoDb:                   repoDb,
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
	return s.repoDb.Exists(ctx, login)
}

func (s *Service) Register(ctx context.Context, user *model.GetPars) (*model.User, error) {
	var err error
	user.Password, err = s.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	err = s.repoDb.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	u, _, err := s.repoDb.Get(ctx, &model.GetPars{
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

	bonusTransactions, err := s.bonusTransactionsService.List(ctx, &bonus_transactionsModel.ListPars{
		UserID:          &pars.ID,
		TransactionType: bonus_transactionsModel.Debit,
	})

	withdrawn := 0
	for _, v := range bonusTransactions {
		withdrawn += v.Sum
	}

	return &model.UserBalance{
		Current:   user.Balance,
		Withdrawn: withdrawn,
	}, nil
}

func (s *Service) List(ctx context.Context, pars *model.ListPars) ([]*model.User, error) {
	return s.repoDb.List(ctx, pars)
}

func (s *Service) Create(ctx context.Context, obj *model.GetPars) error {
	return s.repoDb.Create(ctx, obj)
}

func (s *Service) Get(ctx context.Context, pars *model.GetPars) (*model.User, bool, error) {
	return s.repoDb.Get(ctx, pars)
}

func (s *Service) Update(ctx context.Context, pars *model.GetPars) error {
	return s.repoDb.Update(ctx, pars)
}

func (s *Service) Delete(ctx context.Context, pars *model.GetPars) error {
	return s.repoDb.Delete(ctx, pars)
}

func (s *Service) Exists(ctx context.Context, orderNumber string) (bool, error) {
	return s.repoDb.Exists(ctx, orderNumber)
}
