package service

import (
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/data"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

type IWalletService interface {
	Get(username string) int
}

func NewWalletService() IWalletService {
	return &WalletService{}
}

type WalletService struct {
	data data.IWalletData
}

func (*WalletService) Get(username string) int { return 0 }

func (s *WalletService) Create(username string) *model.Wallet {
	balance := config.Getconf().InitialBalanceAmount
	return s.data.Create(username, balance)
}
