package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// create represents the create command which creates a new expense
func (s Switch) create() func(string) error {
	return func(cmdName string) error {
		createCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		t := setTitleFlag(createCmd, false)
		c := setCurrencyFlag(createCmd, false)
		p := setPriceFlag(createCmd, false)

		if err := s.parseCmd(createCmd); err != nil {
			return err
		}
		if err := s.checkArgs(createCmd, 3); err != nil {
			return err
		}

		err := s.client.Create(t.value, c.value, p.value)
		if err != nil {
			return errors.Wrap(err, "could not create expense")
		}

		fmt.Println("expense created successfully")
		return nil
	}
}
