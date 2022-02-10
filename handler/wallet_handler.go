package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/service"
)

type IWalletHandler interface {
	HandleWalletEndpoints(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}

func NewWalletHandler(service service.IWalletService) IWalletHandler {
	return &WalletHandler{service: service}
}

type WalletHandler struct {
	service service.IWalletService
}

func (h *WalletHandler) HandleWalletEndpoints(w http.ResponseWriter, r *http.Request) {
	// redirect to related methods according to URI
	rex := regexp.MustCompile(`/[a-z]*`)
	if !rex.MatchString(r.RequestURI) {
		http.Error(w, "URI can't be compiled", http.StatusBadRequest)
	}

	if r.RequestURI == "/" && r.Method == "GET" {
		h.GetAll(w, r)
		return
	}

	if r.Method == "GET" {
		h.Get(w, r)
	} else if r.Method == "PUT" {
		h.Create(w, r)
	} else if r.Method == "POST" {
		h.Update(w, r)
	}
}

func (h *WalletHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	wallets := h.service.GetAllWallets()
	walletBytes, err := json.Marshal(wallets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(walletBytes)
}

func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	username := r.RequestURI[1:]
	balance, err := h.service.GetBalance(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s's wallet balance: %d", username, balance)))
}

func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HANDLER CREATE")
	username := r.RequestURI[1:]
	if len(username) == 0 {
		http.Error(w, "Username can't be empty", http.StatusBadRequest)
		return
	}

	wallet := h.service.CreateWallet(username)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Wallet owner: %s, Wallet balance: %d", wallet.Username, wallet.Balance)))
}

var wallet model.Wallet

func (h *WalletHandler) Update(w http.ResponseWriter, r *http.Request) {
	username := r.RequestURI[1:]
	err := json.NewDecoder(r.Body).Decode(&wallet)
	fmt.Println(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedBalance, err := h.service.UpdateBalance(username, wallet.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Updated wallet balance: %d", updatedBalance)))
}
