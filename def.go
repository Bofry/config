package config

type UnmarshalFunc func(buffer []byte, target interface{}) error
