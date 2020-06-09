package main

// Now that cli.go has processed the flags, validate there are no
// conflicts and drive to the business logic.

import (
	"fmt"
	"log"
	"os"

	"github.com/StackExchange/blackbox/v2/pkg/box"
	"github.com/urfave/cli/v2"
)

var logErr *log.Logger

func init() {
	if logErr == nil {
		logErr = log.New(os.Stderr, "", 0)
	}
}

func allOrSomeFiles(c *cli.Context) error {
	if c.Bool("all") && c.Args().Present() {
		return fmt.Errorf("Can not specify filenames and --all")
	}
	if (!c.Args().Present()) && (!c.Bool("all")) {
		return fmt.Errorf("Must specify at least one file name or --all")
	}
	return nil
}

// Keep these functions in alphabetical order.

func cmdAdminAdd(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf(
			"Must specify at least one admin's GnuPG user-id (i.e. email address)")
	}
	bx := box.NewFromFlags(c)
	return bx.AdminAdd(c.Args().Slice())
}

func cmdAdminList(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	return bx.AdminList()
}

func cmdAdminRemove(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one admin's GnuPG user-id (i.e. email address)")
	}
	bx := box.NewFromFlags(c)
	return bx.AdminRemove(c.Args().Slice())
}

func cmdCat(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	return bx.Cat(c.Args().Slice())
}

func cmdDecrypt(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}

	// The default for --agentcheck is off normally, and on when using --all.
	pauseNeeded := c.Bool("all")
	// If the user used the flag, abide by it.
	if c.IsSet("agentcheck") {
		pauseNeeded = c.Bool("agentcheck")
	}

	bx := box.NewFromFlags(c)
	return bx.Decrypt(c.Args().Slice(),
		c.Bool("overwrite"),
		pauseNeeded,
		c.String("group"),
	)
}

func cmdDiff(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}
	bx := box.NewFromFlags(c)
	return bx.Diff(c.Args().Slice())
}

func cmdEdit(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	return bx.Edit(c.Args().Slice())
}

func cmdEncrypt(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}
	bulk := false
	if c.Bool("all") {
		bulk = c.Bool("bulk") // Only applies to --all
	}
	bx := box.NewFromFlags(c)
	return bx.Encrypt(c.Args().Slice(), bulk, c.String("group"), c.Bool("overwrite"))
}

func cmdFileAdd(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	return bx.FileAdd(c.Args().Slice(), c.Bool("overwrite"))
}

func cmdFileList(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	return bx.FileList()
}

func cmdFileRemove(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	return bx.FileRemove(c.Args().Slice())
}

func cmdInfo(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	return bx.Info()
}

func cmdInit(c *cli.Context) error {
	if c.Args().Len() > 1 {
		return fmt.Errorf("This command takes one or two arguments")
	}
	bx := box.NewUninitialized()
	return bx.Init(c.Args().First(), c.String("vcs"))
}

func cmdReencrypt(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}
	bx := box.NewFromFlags(c)
	return bx.Reencrypt(c.Args().Slice())
}

func cmdShred(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}
	bx := box.NewFromFlags(c)
	return bx.Shred(c.Args().Slice())
}

func cmdStatus(c *cli.Context) error {
	if c.Bool("all") && c.Args().Present() {
		return fmt.Errorf("Can not specify filenames and --all")
	}
	bx := box.NewFromFlags(c)
	return bx.Status(c.Args().Slice(), c.Bool("name-only"), c.String("type"))
}

// These are "secret" commands used by the integration tests.

func testingInit(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("No args required")
	}
	fmt.Printf(
		"c.String(vcs) reports %q\n",
		c.String("vcs"),
	)
	bx := box.NewBare(c.String("vcs"))
	return bx.TestingInitRepo()
}
