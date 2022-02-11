package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/service"
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
	GetAll() map[string]*model.Wallet
	Get(username string) (*model.Wallet, error)
	Create(username string) *model.Wallet
	Update(username string, balance int) (*model.Wallet, error)
}

func NewWallet(service WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

// WalletHandler uses the WalletService's interface
// which was created in the handler (consumer)
type WalletHandler struct {
	service WalletService
}

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	httpErr := HTTPError{
		Code:    code,
		Message: msg,
	}
	bytesResponse, _ := json.Marshal(httpErr)
	w.WriteHeader(code)
	w.Write(bytesResponse)
}

func (h *WalletHandler) HandleWalletEndpoints(w http.ResponseWriter, r *http.Request) {
	// Clean the last forward slash if there is one
	r.RequestURI = strings.TrimSuffix(r.RequestURI, "/")
	pathsAndParams := strings.Split(r.RequestURI, forwardSlash)

	// Set content type to json
	w.Header().Set("Content-Type", "application/json")

	// If first element of pathsAndParams is not an empty string,
	// this means the URI does not start with a forward slash,
	// so we raise a bad request error
	// We also raise a bad request error if client tries to write
	// more than two forward slashes because we do not support
	// that sort of endpoint
	if len(pathsAndParams[0]) != 0 || len(pathsAndParams) > 2 {
		writeErr(w, http.StatusNotFound, "Invalid URI")
		return
	}

	// If URI was just forward slash, now it's empty because we trimmed it
	if r.RequestURI == "" {
		if r.Method == "GET" {
			h.GetAll(w, r)
		} else {
			writeErr(w, http.StatusNotFound, "Invalid endpoint")
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
	wallets := h.service.GetAll()

	response, err := json.Marshal(wallets)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *WalletHandler) Get(w http.ResponseWriter, r *http.Request) {
	username := route.PullParam("username")
	wallet, err := h.service.Get(username)
	if err != nil && errors.Is(err, service.ErrWalletNotExists) {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse, err := json.Marshal(wallet)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *WalletHandler) Create(w http.ResponseWriter, r *http.Request) {
	username := route.PullParam("username")
	wallet := h.service.Create(username)

	jsonResponse, err := json.Marshal(wallet)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

var wallet model.Wallet

func (h *WalletHandler) Update(w http.ResponseWriter, r *http.Request) {
	username := route.PullParam("username")
	err := json.NewDecoder(r.Body).Decode(&wallet)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedWallet, err := h.service.Update(username, wallet.Balance)
	if err != nil && errors.Is(err, service.ErrWalletNotExists) {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil && errors.Is(err, service.ErrBalanceBelowLimit) {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	} else if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse, err := json.Marshal(updatedWallet)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
