package repository

import (
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

func NewWallet(wallets map[string]*model.Wallet) *DefaultWalletRepository {
	return &DefaultWalletRepository{wallets: wallets}
}

type DefaultWalletRepository struct {
	wallets map[string]*model.Wallet
}

func (r *DefaultWalletRepository) Exists(username string) bool {
	_, ok := r.wallets[username]
	return ok
}

func (r *DefaultWalletRepository) GetAll() map[string]*model.Wallet {
	return r.wallets
}

func (r *DefaultWalletRepository) Get(username string) *model.Wallet {
	return r.wallets[username]
}

func (r *DefaultWalletRepository) Create(username string, initialBalance int) *model.Wallet {
	wallet := &model.Wallet{
		Username: username,
		Balance:  initialBalance,
	}
	r.wallets[username] = wallet
	return wallet
}

func (r *DefaultWalletRepository) Update(username string, balance int) *model.Wallet {
	r.wallets[username].Balance += balance
	return r.wallets[username]
}
