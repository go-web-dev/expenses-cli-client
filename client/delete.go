package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// delete represents the delete command which delete one expense by a given id
func (s Switch) delete() func(string) error {
	return func(cmdName string) error {
		deleteCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		ids := setIDsFlag(deleteCmd)

		if err := s.parseCmd(deleteCmd); err != nil {
			return err
		}
		if err := s.checkArgs(deleteCmd, 1); err != nil {
			return err
		}
		if len(ids.value) == 0 {
			return errors.New("id of the expense must be provided")
		}

		err := s.client.Delete(ids.value[0])
		if err != nil {
			return errors.Wrap(err, "could not delete expense")
		}

		fmt.Printf("expense with id: %s deleted successfully\n", ids.value[0])
		return nil
	}
}
