package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// logout represents the logout command which logs the user out and removes the access token from file
func (s Switch) logout() func(string) error {
	return func(cmdName string) error {
		logoutCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		if err := s.parseCmd(logoutCmd); err != nil {
			return err
		}

		err := s.client.Logout()
		if err != nil {
			return errors.Wrap(err, "could not log out the user")
		}
		if err := saveCredentials(Credentials{}); err != nil {
			return errors.Wrap(err, "could not clear credentials from file")
		}

		fmt.Println("successfully logged out the user")
		return nil
	}
}
