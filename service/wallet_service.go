package service

import (
	"fmt"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

type WalletRepository interface {
	GetAllWallets() map[string]int
	GetBalance(username string) (int, error)
	CreateWallet(username string, balance int) *model.Wallet
	UpdateBalance(username string, balance int, minimumBalance int) (int, error)
}

func NewWallet(repo WalletRepository) *DefaultWalletService {
	return &DefaultWalletService{
		repository: repo,
	}
}

type DefaultWalletService struct {
	repository WalletRepository
}

func (s *DefaultWalletService) GetAllWallets() map[string]int {
	return s.repository.GetAllWallets()
}

func (s *DefaultWalletService) GetBalance(username string) (int, error) {
	return s.repository.GetBalance(username)
}

func (s *DefaultWalletService) CreateWallet(username string) *model.Wallet {
	fmt.Println("SERVICE CREATE")
	balance := config.Getconf().InitialBalanceAmount
	fmt.Println(balance)
	wallet := s.repository.CreateWallet(username, balance)
	fmt.Println(wallet)
	return wallet
}

func (s *DefaultWalletService) UpdateBalance(username string, balance int) (int, error) {
	minimumBalance := config.Getconf().MinimumBalanceAmount
	return s.repository.UpdateBalance(username, balance, minimumBalance)
}
