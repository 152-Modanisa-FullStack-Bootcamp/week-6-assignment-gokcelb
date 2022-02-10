package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
)

const forwardSlash = "/"

// Route is where I will store my route parameters
type Route struct {
	params map[string]string
}

func (route *Route) PullParam(param string) string {
	if val, ok := route.params[param]; ok {
		return val
	}
	return ""
}

func (route *Route) PushParam(param string, paramVal string) {
	route.params[param] = paramVal
}

var route = &Route{
	params: make(map[string]string),
}

// Consumer-side interface
type WalletService interface {
	GetAllWallets() map[string]int
	GetBalance(username string) (int, error)
	CreateWallet(username string) *model.Wallet
	UpdateBalance(username string, balance int) (int, error)
}

func NewWallet(service WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

// WalletHandler uses the WalletService's interface
// which was created in the handler (consumer)
type WalletHandler struct {
	service WalletService
}

func (h *WalletHandler) HandleWalletEndpoints(w http.ResponseWriter, r *http.Request) {
	// Clean the last forward slash if there is one
	r.RequestURI = strings.TrimSuffix(r.RequestURI, "/")
	pathsAndParams := strings.Split(r.RequestURI, forwardSlash)

	// If first element of pathsAndParams is not an empty string,
	// this means the URI does not start with a forward slash,
	// so we raise a bad request error
	// We also raise a bad request error if client tries to write
	// more than two forward slashes because we do not support
	// that sort of endpoint
	if len(pathsAndParams[0]) != 0 || len(pathsAndParams) > 2 {
		http.Error(w, "Invalid URI", http.StatusNotFound)
		return
	}

	// If URI was just forward slash, now it's empty because we trimmed it
	if r.RequestURI == "" {
		if r.Method == "GET" {
			h.GetAll(w, r)
		} else {
			http.Error(w, "Invalid endpoint", http.StatusNotFound)
		}
		return
	}

	// If the method progresses until this point, it means the endpoint
	// contains a route parameter. I assume that everything after the
	// forward slash comprises the username, and push it as a route
	// parameter, then I will pull this username route parameter
	username := pathsAndParams[1]
	route.PushParam("username", username)

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

	response, err := json.Marshal(wallets)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	username := route.PullParam("username")
	if len(username) == 0 {
		http.Error(w, "Route parameter error", http.StatusInternalServerError)
		return
	}

	balance, err := h.service.GetBalance(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{}
	response["username"] = username
	response["balance"] = balance
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	username := route.PullParam("username")
	if len(username) == 0 {
		http.Error(w, "Username cannot be empty", http.StatusBadRequest)
		return
	}

	wallet := h.service.CreateWallet(username)

	jsonResponse, err := json.Marshal(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

var wallet model.Wallet

func (h *WalletHandler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: Should return an error if wallet doesn't exist
	username := route.PullParam("username")
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedBalance, err := h.service.UpdateBalance(username, wallet.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{}
	response["username"] = username
	response["balance"] = updatedBalance
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
