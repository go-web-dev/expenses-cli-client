package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// update represents the update command which updates an existing expense
func (s Switch) update() func(string) error {
	return func(cmdName string) error {
		updateCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		t := setTitleFlag(updateCmd, true)
		c := setCurrencyFlag(updateCmd, true)
		p := setPriceFlag(updateCmd, true)
		ids := setIDsFlag(updateCmd)

		if err := s.parseCmd(updateCmd); err != nil {
			return err
		}
		if err := s.checkArgs(updateCmd, 2); err != nil {
			return err
		}
		if len(ids.value) == 0 {
			return errors.New("id of the expense must be provided")
		}

		err := s.client.Update(ids.value[0], t.value, c.value, p.value)
		if err != nil {
			return errors.Wrap(err, "could not update expense")
		}

		fmt.Println("expense updated successfully")
		return nil
	}
}
