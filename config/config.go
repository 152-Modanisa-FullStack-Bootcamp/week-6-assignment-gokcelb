package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

const (
	defaultPath   = ".config\\"
	fileExtention = ".json"
	rootDir       = "c:\\users\\bilgi\\desktop\\bc\\week-6-assignment-gokcelb"
)

var c = &Conf{}

type Conf struct {
	InitialBalanceAmount int `json:"initialBalanceAmount"`
	MinimumBalanceAmount int `json:"minimumBalanceAmount"`
}

func Getconf() *Conf {
	log.Print("CONFIG INITIALIZED")
	// We will change path according to current working directory
	// because while the application normally runs on the root directory,
	// it runs on the service directory when we run tests
	path := defaultPath
	currWD, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Print("current directory:", currWD)

	if strings.ToLower(currWD) != rootDir {
		path = "..\\.config\\"
	}
	log.Print("path:", path)

	file, err := os.Open(path + Getenv() + fileExtention)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	read, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(read, c); err != nil {
		panic(err)
	}

	return c
}
