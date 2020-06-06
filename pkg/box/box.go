package box

// box implements the business logic for all operations on the blackbox.

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type box struct {
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

func (bx *box) AdminAdd([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminAdd")
}

func (bx *box) AdminList() error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminList")
}

func (bx *box) AdminRemove([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminRemove")
}

func (bx *box) Cat([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Cat")
}

func (bx *box) Decrypt(names []string, overwrite bool, bulk bool, setgroup string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Decrypt")
}

func (bx *box) Diff([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Diff")
}

func (bx *box) Edit([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Edit")
}

func (bx *box) Encrypt(names []string, bulk bool, setgroup string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Encrypt")
}

func (bx *box) FileAdd(names []string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileAdd")
}

func (bx *box) FileList() error {
	return fmt.Errorf("NOT IMPLEMENTED: FileList")
}

func (bx *box) FileRemove(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileRemove")
}

func (bx *box) Info() error {
	return fmt.Errorf("NOT IMPLEMENTED: Info")
}

func (bx *box) Init() error {
	return fmt.Errorf("NOT IMPLEMENTED: Init")
}

func (bx *box) Reencrypt(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Reencrypt")
}

func (bx *box) Shred(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Shred")
}

func (bx *box) Status(names []string, mode StatusMode, nameOnly bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Status")
}
