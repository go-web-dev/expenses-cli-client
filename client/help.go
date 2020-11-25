package client

import (
	"fmt"
	"os"
	"strings"
)

// Help prints a useful message about command usage
func (s Switch) Help() {
	var help string
	for name := range s.commands {
		help += name + strings.Repeat(" ", 10-len(name)) + "\t --help\n"
	}
	fmt.Printf("Usage of %s:\n<command> [<args>]\n%s", os.Args[0], help)
}
