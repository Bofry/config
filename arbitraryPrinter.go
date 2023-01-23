package config

import (
	"io"
)

var _ Printer = new(ArbitraryPrinter)

type ArbitraryPrinter struct {
	writer    io.Writer
	successor Printer
}

func NewArbitraryPrinter(writer io.Writer, successor Printer) *ArbitraryPrinter {
	return &ArbitraryPrinter{
		writer:    writer,
		successor: successor,
	}
}

func (p *ArbitraryPrinter) Print(target interface{}) error {
	if v, ok := target.(Printable); ok {
		err := v.Output(p.writer)
		return err
	}
	return p.successor.Print(target)
}
