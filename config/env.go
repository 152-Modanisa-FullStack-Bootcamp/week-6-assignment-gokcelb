package config

import "os"

const key = "APP_ENV"
const defaultEnv = "local"

func Getenv() string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return defaultEnv
}
