package client

import (
	"flag"
	"fmt"

	"github.com/pkg/errors"
)

// signup represents the signup command which signs the user up and saves the access token to file
func (s Switch) signup() func(string) error {
	return func(cmdName string) error {
		signupCmd := flag.NewFlagSet(cmdName, flag.ExitOnError)
		email, pwd := setEmailFlag(signupCmd), setPasswordFlag(signupCmd)
		if err := s.parseCmd(signupCmd); err != nil {
			return err
		}
		if err := s.checkArgs(signupCmd, 2); err != nil {
			return err
		}

		_, err := s.client.Signup(email.value, pwd.value)
		if err != nil {
			return errors.Wrap(err, "could not sign up the user")
		}

		credentials := Credentials{AccessToken: "some-access-token"}
		err = saveCredentials(credentials)
		if err != nil {
			return errors.Wrap(err, "could not save credentials to file")
		}

		fmt.Println("successfully signed up the user")
		return nil
	}
}
