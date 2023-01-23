package config

import (
	"fmt"
	"io"
)

var _ Printer = new(CommonPrinter)

type CommonPrinter struct {
	writer io.Writer
}

func NewCommonPrinter(writer io.Writer) *CommonPrinter {
	return &CommonPrinter{
		writer: writer,
	}
}

func (p *CommonPrinter) Print(v interface{}) error {
	fmt.Fprintf(p.writer, "%+v\n", v)
	return nil
}
