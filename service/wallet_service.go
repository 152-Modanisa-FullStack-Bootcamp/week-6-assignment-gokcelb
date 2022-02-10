package service

import (
	"fmt"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/data"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

type IWalletService interface {
	GetAllWallets() map[string]int
	GetBalance(username string) (int, error)
	CreateWallet(username string) *model.Wallet
	UpdateBalance(username string, balance int) (int, error)
}

func NewWalletService(data data.IWalletData) IWalletService {
	return &WalletService{
		data: data,
	}
}

type WalletService struct {
	data data.IWalletData
}

func (s *WalletService) GetAllWallets() map[string]int {
	return s.data.GetAllWallets()
}

func (s *WalletService) GetBalance(username string) (int, error) {
	return s.data.GetBalance(username)
}

func (s *WalletService) CreateWallet(username string) *model.Wallet {
	fmt.Println("SERVICE CREATE")
	balance := config.Getconf().InitialBalanceAmount
	fmt.Println(balance)
	wallet := s.data.CreateWallet(username, balance)
	fmt.Println(wallet)
	return wallet
}

func (s *WalletService) UpdateBalance(username string, balance int) (int, error) {
	minimumBalance := config.Getconf().MinimumBalanceAmount
	return s.data.UpdateBalance(username, balance, minimumBalance)
}
