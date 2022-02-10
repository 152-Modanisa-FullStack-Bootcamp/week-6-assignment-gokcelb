package data

import (
	"errors"
	"fmt"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

type IWalletData interface {
	GetAllWallets() map[string]int
	GetBalance(username string) (int, error)
	CreateWallet(username string, balance int) *model.Wallet
	UpdateBalance(username string, balance int, minimumBalance int) (int, error)
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

func (d *WalletData) GetAllWallets() map[string]int {
	return d.wallets
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

func (d *WalletData) UpdateBalance(username string, balance int, minimumBalance int) (int, error) {
	newBalance := d.wallets[username] + balance
	if newBalance < minimumBalance {
		return -1, errors.New(fmt.Sprintf("Wallet balance cannot be lower than %d", minimumBalance))
	}
	d.wallets[username] = newBalance
	return newBalance, nil
}
