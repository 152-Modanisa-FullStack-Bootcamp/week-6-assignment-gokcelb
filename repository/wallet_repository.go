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

func (r *DefaultWalletRepository) GetAll() []model.Wallet {
	wallets := []model.Wallet{}
	for _, wallet := range r.wallets {
		wallets = append(wallets, *wallet)
	}
	return wallets
}

func (r *DefaultWalletRepository) Get(username string) model.Wallet {
	if wallet, ok := r.wallets[username]; ok {
		return *wallet
	}
	return model.Wallet{}
}

func (r *DefaultWalletRepository) Save(wallet *model.Wallet) {
	r.wallets[wallet.Username] = wallet
}

func (r *DefaultWalletRepository) Update(username string, balance int) model.Wallet {
	r.wallets[username].Balance = balance
	return *r.wallets[username]
}

func (r *DefaultWalletRepository) Delete(username string) {
	delete(r.wallets, username)
}
