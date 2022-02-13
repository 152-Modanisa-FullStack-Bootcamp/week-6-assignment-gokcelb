package config

import (
	"encoding/json"
	"os"
)

type Conf struct {
	InitialBalanceAmount int `json:"initialBalanceAmount"`
	MinimumBalanceAmount int `json:"minimumBalanceAmount"`
}

func Read(path string) (*Conf, error) {
	contentBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Conf
	err = json.Unmarshal(contentBytes, &c)
	return &c, err
}
