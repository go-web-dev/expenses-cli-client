package client

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// BackendHTTPClient represents the HTTP client for communicating with the Backend API
type BackendHTTPClient interface {
	GetAll(page, pageSize string) ([]byte, error)
	GetByIDs(ids ...string) ([]byte, error)
	Create(title, currency string, price float64) error
	Update(id, title, currency string, price float64) error
	Delete(id string) error
	Login(email, password string) ([]byte, error)
	Logout() error
	Signup(email, password string) ([]byte, error)
}

// NewSwitch creates a new instance of command Switch
func NewSwitch(uri string) Switch {
	httpClient := NewHTTPClient(uri)
	s := Switch{client: httpClient, backendAPIURL: uri}
	s.commands = map[string]func() func(string) error{
		"get-all":    s.getAll,
		"get-by-ids": s.getByIDs,
		"create":     s.create,
		"update":     s.update,
		"delete":     s.delete,
		"login":      s.login,
		"logout":     s.logout,
		"signup":     s.signup,
	}
	return s
}

// Switch represents CLI command switch
type Switch struct {
	client        BackendHTTPClient
	backendAPIURL string
	commands      map[string]func() func(string) error
}

// Switch analyses the CLI args and executes the given command
func (s Switch) Switch() error {
	cmdName := os.Args[1]
	cmd, ok := s.commands[os.Args[1]]
	if !ok {
		return fmt.Errorf("invalid command '%s'", cmdName)
	}
	return cmd()(cmdName)
}

// parseCmd parses sub-command flags
func (s Switch) parseCmd(cmd *flag.FlagSet) error {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		return errors.Wrap(err, "could not parse '"+cmd.Name()+"' flags")
	}
	return nil
}

// checkArgs checks if the number of passed args for a command is greater or equal to min args
func (s Switch) checkArgs(cmd *flag.FlagSet, minArgs int) error {
	if cmd.NFlag() < minArgs {
		fmt.Printf(
			"incorect use of %s\n%s %s --help\n",
			os.Args[1], os.Args[0], os.Args[1],
		)
		return fmt.Errorf(
			"%s expects at least: %d arg(s), %d provided",
			cmd.Name(),
			minArgs,
			cmd.NFlag(),
		)
	}
	return nil
}
