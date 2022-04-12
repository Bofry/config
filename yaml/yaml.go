package yaml

import "github.com/Bofry/config/internal/yaml"

func LoadFile(filepath string, target interface{}) error {
	return yaml.LoadFile(filepath, target)
}

func LoadBytes(buffer []byte, target interface{}) error {
	return yaml.LoadBytes(buffer, target)
}
