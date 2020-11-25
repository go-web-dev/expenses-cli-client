package client

import (
	"encoding/json"
	"io"
	"os"

	"github.com/pkg/errors"
)

const credentialsFileName = ".credentials.json"

// Credentials represents the user credentials after successful login/signup
type Credentials struct {
	AccessToken string `json:"access_token"`
}

func saveCredentials(credentials Credentials) error {
	_, err := os.Stat(credentialsFileName)
	var credentialsFile *os.File
	switch err.(type) {
	case *os.PathError:
		f, err := os.Create(credentialsFileName)
		if err != nil {
			return errors.Wrap(err, "could not create credentials file")
		}
		credentialsFile = f
	case nil:
		f, err := os.OpenFile(credentialsFileName, os.O_RDWR, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "could not open credentials file")
		}
		credentialsFile = f
	}
	err = json.NewEncoder(credentialsFile).Encode(credentials)
	if err != nil {
		return errors.Wrap(err, "could not encode json credentials")
	}
	return nil
}

func readCredentials() (Credentials, error) {
	var credentials Credentials
	credentialsFile, err := os.Open(credentialsFileName)
	if err != nil {
		return Credentials{}, errors.Wrap(err, "could not open credentials file")
	}
	err = json.NewDecoder(credentialsFile).Decode(&credentials)
	if err != nil && err != io.EOF {
		return Credentials{}, errors.Wrap(err, "could not decode json credentials")
	}
	return credentials, nil
}
