package box

// box implements the box model.

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type box struct {
	Admins []string
}

func NewFromFlags(c *cli.Context) *box {
	return &box{}
}

type StatusMode int

const (
	Itemized StatusMode = iota
	All
	Unchanged
	Changed
)

func (bx *box) getAdmins() ([]string, error) {
	if len(bx.Admins) != 0 {
		return bx.Admins, nil
	}

	fmt.Printf("Would load\n")
	bx.Admins = nil

	return bx.Admins, nil
}
