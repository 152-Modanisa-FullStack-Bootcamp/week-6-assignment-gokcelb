package data

import (
	"errors"
	"fmt"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

type IWalletData interface {
	GetBalance(username string) (int, error)
	CreateWallet(username string, balance int) *model.Wallet
}

// initialize WalletData with empty map,
// otherwise map values are nil, instead of zero value
var wallets = make(map[string]int)

func NewWalletData() IWalletData {
	return &WalletData{wallets: wallets}
}

type WalletData struct {
	wallets map[string]int
}

func (d *WalletData) GetBalance(username string) (int, error) {
	if val, ok := d.wallets[username]; ok {
		return val, nil
	}
	return -1, errors.New("No wallet belonging to the given username exists")
}

func (d *WalletData) CreateWallet(username string, balance int) *model.Wallet {
	fmt.Println("DATA CREATE")
	// if there is already a wallet registered under
	// the given name, send the already existing balance value
	// otherwise, create wallet and send the given balance value
	if val, ok := d.wallets[username]; ok {
		return &model.Wallet{
			Username: username,
			Balance:  val,
		}
	}
	d.wallets[username] = balance
	return &model.Wallet{
		Username: username,
		Balance:  balance,
	}
}
