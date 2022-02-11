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
	r := repository.NewWallet(wallets)
	s := service.NewWallet(r)
	h := handler.NewWallet(s)

	http.HandleFunc("/", h.HandleWalletEndpoints)
	http.ListenAndServe(":8000", nil)
}
