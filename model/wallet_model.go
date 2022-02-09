package model

type Wallet struct {
	Username string
	Balance  int `json:"balance"`
}
