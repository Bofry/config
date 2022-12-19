package env

import "github.com/Bofry/config/internal/env"

func Process(prefix string, target interface{}) error {
	return env.Process(prefix, target)
}

func LoadDotEnv(target interface{}) error {
	return env.LoadDotEnv(target)
}

func LoadDotEnvFile(filepath string, target interface{}) error {
	return env.LoadDotEnvFile(filepath, target)
}
