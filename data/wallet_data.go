package data

import (
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

type IWalletData interface {
	Create(username string, balance int) *model.Wallet
}

func New() IWalletData {
	return &WalletData{}
}

type WalletData struct {
	wallets map[string]int
}

func (d *WalletData) Create(username string, balance int) *model.Wallet {
	if val, ok := d.wallets[username]; ok {
		return &model.Wallet{
			Username: username,
			Balance:  val,
		}
	} else {
		return &model.Wallet{
			Username: username,
			Balance:  balance,
		}
	}
}
