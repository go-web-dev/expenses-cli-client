package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// login represents the login command which logs the user in and saves the access token to file
func (s Switch) login() func(string) error {
	return func(cmdName string) error {
		loginCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		email, pwd := setEmailFlag(loginCmd), setPasswordFlag(loginCmd)
		if err := s.parseCmd(loginCmd); err != nil {
			return err
		}
		if err := s.checkArgs(loginCmd, 2); err != nil {
			return err
		}

		_, err := s.client.Login(email.value, pwd.value)
		if err != nil {
			return errors.Wrap(err, "could not login user")
		}

		credentials := Credentials{AccessToken: "some-access-token"}
		err = saveCredentials(credentials)
		if err != nil {
			return errors.Wrap(err, "could not save credentials to file")
		}

		fmt.Println("successfully logged in")
		return nil
	}
}
