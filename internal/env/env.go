package env

import (
	"os"
	"strings"

	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/valuebinder"
	"github.com/joho/godotenv"
)

const (
	TagName = "env"
)

func Process(prefix string, target interface{}) error {
	if len(prefix) > 0 {
		prefix += "_"
	}

	prototype, err := structproto.Prototypify(target, &structproto.StructProtoResolveOption{
		TagName: TagName,
	})
	if err != nil {
		return err
	}

	var table structproto.FieldValueMap = make(structproto.FieldValueMap)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		name, value := parts[0], parts[1]
		if strings.HasPrefix(name, prefix) {
			table[name[len(prefix):]] = value
		}
	}
	err = prototype.BindIterator(table, valuebinder.BuildStringBinder)
	if err != nil {
		return err
	}
	return nil
}

func LoadDotEnv(target interface{}) error {
	var err error
	err = godotenv.Load()
	if err != nil {
		return err
	}

	return Process("", target)
}

func LoadDotEnvFile(filepath string, target interface{}) error {
	var err error
	path := os.ExpandEnv(filepath)
	err = godotenv.Load(path)
	if err != nil {
		return err
	}

	return Process("", target)
}
