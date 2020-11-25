package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// getByIDs represents the get-by-ids command which fetches all expenses by a list of given ids
func (s Switch) getByIDs() func(string) error {
	return func(cmdName string) error {
		getByIDsCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		ids := setIDsFlag(getByIDsCmd)
		if err := s.parseCmd(getByIDsCmd); err != nil {
			return err
		}
		if err := s.checkArgs(getByIDsCmd, 1); err != nil {
			return err
		}
		if len(ids.value) == 0 {
			return errors.New("at least one expense id must be provided")
		}

		res, err := s.client.GetByIDs(ids.value...)
		if err != nil {
			return errors.Wrap(err, "could not fetch expenses")
		}

		fmt.Printf("expenses fetched successfully:\n%s\n", string(res))
		return nil
	}
}
