package config

var _ Printer = NonePrinter{}

type NonePrinter struct{}

func (NonePrinter) Print(target interface{}) error {
	return nil
}
