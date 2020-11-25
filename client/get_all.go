package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// getAll represents the get-all command which fetches all expenses with pagination
func (s Switch) getAll() func(string) error {
	return func(cmdName string) error {
		getAllCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		page, pageSize := setPageFlag(getAllCmd), setPageSizeFlag(getAllCmd)
		if err := s.parseCmd(getAllCmd); err != nil {
			return err
		}

		res, err := s.client.GetAll(page.value, pageSize.value)
		if err != nil {
			return errors.Wrap(err, "could not fetch expenses")
		}

		fmt.Printf("expenses fetched successfully:\n%s\n", string(res))
		return nil
	}
}
