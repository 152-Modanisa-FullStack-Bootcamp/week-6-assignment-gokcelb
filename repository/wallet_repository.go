package repository

import (
	"errors"
	"fmt"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

// initialize WalletData with empty map,
// otherwise map values are nil, instead of zero value
var wallets = make(map[string]int)

func NewWallet() *DefaultWalletRepository {
	return &DefaultWalletRepository{wallets: wallets}
}

type DefaultWalletRepository struct {
	wallets map[string]int
}

func (r *DefaultWalletRepository) GetAllWallets() map[string]int {
	return r.wallets
}

func (r *DefaultWalletRepository) GetBalance(username string) (int, error) {
	if val, ok := r.wallets[username]; ok {
		return val, nil
	}
	return -1, errors.New("No wallet belonging to the given username exists")
}

func (r *DefaultWalletRepository) CreateWallet(username string, balance int) *model.Wallet {
	fmt.Println("DATA CREATE")
	// if there is already a wallet registered under
	// the given name, send the already existing balance value
	// otherwise, create wallet and send the given balance value
	if val, ok := r.wallets[username]; ok {
		return &model.Wallet{
			Username: username,
			Balance:  val,
		}
	}
	r.wallets[username] = balance
	return &model.Wallet{
		Username: username,
		Balance:  balance,
	}
}

func (r *DefaultWalletRepository) UpdateBalance(username string, balance, minimumBalance int) (int, error) {
	newBalance := r.wallets[username] + balance
	if newBalance < minimumBalance {
		return -1, errors.New(fmt.Sprintf("Wallet balance cannot be lower than %d", minimumBalance))
	}
	r.wallets[username] = newBalance
	return newBalance, nil
}
