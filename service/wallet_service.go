package service

import (
	"fmt"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

type WalletRepository interface {
	Exists(username string) bool
	GetAll() map[string]*model.Wallet
	Get(username string) *model.Wallet
	Create(username string, balance int) *model.Wallet
	Update(username string, balance int) *model.Wallet
}

func NewWallet(repo WalletRepository) *DefaultWalletService {
	return &DefaultWalletService{
		repository: repo,
	}
}

type DefaultWalletService struct {
	repository WalletRepository
}

func (s *DefaultWalletService) GetAll() map[string]*model.Wallet {
	return s.repository.GetAll()
}

func (s *DefaultWalletService) Get(username string) (*model.Wallet, error) {
	if !s.repository.Exists(username) {
		return nil, fmt.Errorf("No wallet belonging to %s exists", username)
	}

	return s.repository.Get(username), nil
}

func (s *DefaultWalletService) Create(username string) *model.Wallet {
	if s.repository.Exists(username) {
		return s.repository.Get(username)
	}

	// If wallet does not already exists, send config's initial balance amount
	// as balance in order to create new wallet
	initialBalance := config.Getconf().InitialBalanceAmount
	wallet := s.repository.Create(username, initialBalance)
	return wallet
}

func (s *DefaultWalletService) Update(username string, newBalance int) (*model.Wallet, error) {
	// Check if we can successfully get the wallet, if not throw not exists error
	currentWallet, err := s.Get(username)
	if err != nil {
		return nil, err
	}

	// If current balance + new balance is below minimum balance amount, return error
	minimumBalance := config.Getconf().MinimumBalanceAmount
	if currentWallet.Balance+newBalance < minimumBalance {
		return nil, fmt.Errorf("Wallet balance cannot be lower than %d", minimumBalance)
	}

	// Update without any problems, since we handled all errors
	return s.repository.Update(username, newBalance), nil
}
