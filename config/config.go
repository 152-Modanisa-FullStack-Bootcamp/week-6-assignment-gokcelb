package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

const (
	defaultPath   = ".config\\"
	fileExtention = ".json"
	serviceDir    = "c:\\Users\\bilgi\\Desktop\\bc\\week-6-assignment-gokcelb\\service"
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
	var path string
	currWD, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Print("current directory:", currWD)

	if currWD == serviceDir {
		path = "..\\.config\\"
	} else {
		path = defaultPath
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
	err = json.Unmarshal(read, c)
	if err != nil {
		panic(err)
	}

	return c
}
