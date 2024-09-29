package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Bofry/config/internal/env"
	"github.com/Bofry/config/internal/flag"
	"github.com/Bofry/config/internal/json"
	"github.com/Bofry/config/internal/resource"
	"github.com/Bofry/config/internal/yaml"
	"github.com/Bofry/structproto"
)

type ConfigurationService struct {
	target interface{}
}

func NewConfigurationService(target interface{}) *ConfigurationService {
	instance := ConfigurationService{
		target: target,
	}
	return &instance
}

func (service *ConfigurationService) LoadEnvironmentVariables(prefix string) *ConfigurationService {
	err := env.Process(prefix, service.target)
	if err != nil {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadDotEnv() *ConfigurationService {
	err := env.LoadDotEnv(service.target)
	if err != nil && os.IsExist(err) {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadDotEnvFile(filepath string) *ConfigurationService {
	err := env.LoadDotEnvFile(filepath, service.target)
	if err != nil && os.IsExist(err) {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadCommandArguments() *ConfigurationService {
	err := flag.Process(service.target)
	if err != nil {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadJsonFile(filepath string) *ConfigurationService {
	err := json.LoadFile(filepath, service.target)
	if err != nil && os.IsExist(err) {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadJsonBytes(buffer []byte) *ConfigurationService {
	err := json.LoadBytes(buffer, service.target)
	if err != nil {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadYamlFile(filepath string) *ConfigurationService {
	err := yaml.LoadFile(filepath, service.target)
	if err != nil && os.IsExist(err) {
		panic(fmt.Errorf("config: %#v", err))
	}
	return service
}

func (service *ConfigurationService) LoadYamlBytes(buffer []byte) *ConfigurationService {
	err := yaml.LoadBytes(buffer, service.target)
	if err != nil {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadResource(baseDir string) *ConfigurationService {
	err := resource.Process(baseDir, service.target)
	if err != nil {
		panic(fmt.Errorf("config: %v", err))
	}
	return service
}

func (service *ConfigurationService) LoadFile(fullpath string, unmarshal UnmarshalFunc) *ConfigurationService {
	path := os.ExpandEnv(fullpath)
	buffer, err := os.ReadFile(path)
	if err != nil && os.IsExist(err) {
		panic(fmt.Errorf("config: %#v", err))
	}

	err = unmarshal(buffer, service.target)
	if err != nil {
		panic(fmt.Errorf("config: %#v", err))
	}
	return service
}

func (service *ConfigurationService) LoadBytes(buffer []byte, unmarshal UnmarshalFunc) *ConfigurationService {
	err := unmarshal(buffer, service.target)
	if err != nil {
		panic(fmt.Errorf("config: %#v", err))
	}
	return service
}

func (service *ConfigurationService) ExpandEnv(prefix string) error {
	if len(prefix) > 0 {
		prefix += "_"
	}

	err := service.Map(func(field structproto.FieldInfo, rv reflect.Value) error {
		var err error
		defer func() {
			ex := recover()
			if v, ok := ex.(error); ok {
				err = v
			} else {
				err = fmt.Errorf("%+v", err)
			}
		}()

		switch rv.Kind() {
		case reflect.String:
			if !rv.IsZero() {
				val := os.Expand(rv.String(), func(s string) string {
					name := prefix + s
					v := os.Getenv(name)
					if len(v) == 0 {
						panic(fmt.Errorf("missing environment variable '%s'", name))
					}
					return v
				})
				rv.SetString(val)
			}
		}
		return err
	})
	return err
}

func (service *ConfigurationService) Map(mapper structproto.StructMapper) error {
	prototype, err := structproto.Prototypify(service.target,
		&structproto.StructProtoResolveOption{})
	if err != nil {
		return err
	}
	return prototype.Map(mapper)
}

func (service *ConfigurationService) Output() {
	err := __DefaultPrinter.Print(service.target)
	if err != nil {
		panic(fmt.Errorf("config: %#v", err))
	}
}

func (service *ConfigurationService) OutputWithPrinter(printer Printer) {
	err := printer.Print(service.target)
	if err != nil {
		panic(fmt.Errorf("config: %#v", err))
	}
}
