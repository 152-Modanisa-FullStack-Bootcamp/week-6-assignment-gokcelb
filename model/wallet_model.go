package model

type Wallet struct {
	Username string `json:"username,omitempty"`
	Balance  int    `json:"balance"`
}
