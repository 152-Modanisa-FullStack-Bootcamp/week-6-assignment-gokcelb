package service

import (
	"fmt"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

var (
	ErrWalletNotExists   = fmt.Errorf("Wallet does not exist")
	ErrBalanceBelowLimit = fmt.Errorf("Wallet balance cannot be lower than minimum limit")
)

type WalletRepository interface {
	Exists(username string) bool
	GetAll() []model.Wallet
	Get(username string) model.Wallet
	Save(wallet *model.Wallet)
	Update(username string, balance int) model.Wallet
}

func NewWallet(conf *config.Conf, repo WalletRepository) *DefaultWalletService {
	return &DefaultWalletService{
		repo: repo,
		conf: conf,
	}
}

type DefaultWalletService struct {
	repo WalletRepository
	conf *config.Conf
}

func (s *DefaultWalletService) GetAll() []model.Wallet {
	return s.repo.GetAll()
}

func (s *DefaultWalletService) Get(username string) (model.Wallet, error) {
	if !s.repo.Exists(username) {
		return model.Wallet{}, ErrWalletNotExists
	}

	return s.repo.Get(username), nil
}

func (s *DefaultWalletService) Create(username string) model.Wallet {
	if s.repo.Exists(username) {
		return s.repo.Get(username)
	}

	initialBalance := s.conf.InitialBalanceAmount
	wallet := &model.Wallet{Username: username, Balance: initialBalance}
	s.repo.Save(wallet)
	return *wallet
}

func (s *DefaultWalletService) Update(username string, balanceToAdd int) (model.Wallet, error) {
	// Check if we can successfully get the wallet, if not throw not exists error
	currentWallet, err := s.Get(username)
	if err != nil {
		return model.Wallet{}, err
	}

	// If current balance + new balance is below minimum balance amount, return error
	newBalance := currentWallet.Balance + balanceToAdd
	if newBalance < s.conf.MinimumBalanceAmount {
		return model.Wallet{}, ErrBalanceBelowLimit
	}

	// Update without any problems, since we handled all errors
	return s.repo.Update(username, newBalance), nil
}
