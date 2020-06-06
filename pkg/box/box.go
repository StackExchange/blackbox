package box

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Box provides access to a Blackbox.
type Boxer interface {
	AdminAdd([]string) error
	AdminList() error
	AdminRemove([]string) error
}

type box struct {
}

func (bx *box) NewFromFlags(c *cli.Context) error {
}

func (bx *box) AdminAdd([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED")
}

func (bx *box) AdminList() error {
	return fmt.Errorf("NOT IMPLEMENTED")
}

func (bx *box) AdminRemove([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED")
}
