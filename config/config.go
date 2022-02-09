package config

import (
	"encoding/json"
	"io"
	"os"
)

const folderName = ".config/"
const fileExtention = ".json"

type Conf struct {
	InitialBalanceAmount int `json:"initialBalanceAmount"`
	MinimumBalanceAmount int `json:"minimumBalanceAmount"`
}

var c = &Conf{}

func init() {
	file, err := os.Open(folderName + Getenv() + fileExtention)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	read, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(read, c)
	if err != nil {
		panic(err)
	}
}

func Getconf() *Conf {
	return c
}
