package main

import (
	"fmt"
	"net/http"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/handler"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/service"
)

func main() {
	config := config.Getconf()
	fmt.Println(config.InitialBalanceAmount, config.MinimumBalanceAmount)

	s := service.NewWalletService()
	h := handler.NewWalletHandler(s)

	http.HandleFunc("/", h.HandleWalletEndpoints)
	http.ListenAndServe(":8000", nil)
}
