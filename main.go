package main

import (
	"net/http"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/handler"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/repository"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/service"
)

func main() {
	wallets := make(map[string]*model.Wallet)
	walletRepository := repository.NewWallet(wallets)
	walletService := service.NewWallet(walletRepository)
	walletHandler := handler.NewWallet(walletService)

	http.HandleFunc("/", walletHandler.HandleWalletEndpoints)
	http.ListenAndServe(":8000", nil)
}
