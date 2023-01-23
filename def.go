package config

import (
	"io"
	"os"
)

type (
	Printer interface {
		Print(target interface{}) error
	}

	Printable interface {
		Output(writer io.Writer) error
	}
)

type (
	UnmarshalFunc func(buffer []byte, target interface{}) error
)

var (
	__DefaultPrinter = NewArbitraryPrinter(os.Stdout, NewCommonPrinter(os.Stdout))
)
