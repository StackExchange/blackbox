package main

// Now that cli.go has processed the flags, validate there are no
// conflicts and drive to the business logic.

import (
	"fmt"
	"log"
	"os"

	"github.com/StackExchange/blackbox/v2/pkg/bblog"
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

const roError = `This command is disabled due to --config flag being used.
We can not determine if the flag's value is in or out of the repo, and
Blackbox can only work on one repo at a time. If the value is inside the
repo, and you'd like to suggest an algorithm that would let us determine that
automatically, please file a bug. We'd love to have this work better. In the
meanwhile, run this command without the --config flag, perhaps after cd'ing
to the base of the repo.`

// Keep these functions in alphabetical order.

func cmdAdminAdd(c *cli.Context) error {
	if c.NArg() == 0 || c.NArg() > 2 {
		return fmt.Errorf(
			"Must specify one admin's GnuPG user-id (i.e. email address) and optionally the directory of the pubkey data (default ~/.GnuPG)")
	}
	bx := box.NewFromFlags(c)
	if bx.ConfigRO {
		return fmt.Errorf(roError)
	}
	err := bx.AdminAdd(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdAdminList(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	err := bx.AdminList()
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdAdminRemove(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one admin's GnuPG user-id (i.e. email address)")
	}
	bx := box.NewFromFlags(c)
	if bx.ConfigRO {
		return fmt.Errorf(roError)
	}
	err := bx.AdminRemove(c.Args().Slice())
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdCat(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	err := bx.Cat(c.Args().Slice())
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
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
	err := bx.Decrypt(c.Args().Slice(),
		c.Bool("overwrite"),
		pauseNeeded,
		c.String("group"),
	)
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdDiff(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}
	bx := box.NewFromFlags(c)
	err := bx.Diff(c.Args().Slice())
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdEdit(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	err := bx.Edit(c.Args().Slice())
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdEncrypt(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}
	bx := box.NewFromFlags(c)
	err := bx.Encrypt(c.Args().Slice(), c.Bool("shred"))
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdFileAdd(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	if bx.ConfigRO {
		return fmt.Errorf(roError)
	}
	err := bx.FileAdd(c.Args().Slice(), c.Bool("shred"))
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdFileList(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	err := bx.FileList()
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdFileRemove(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("Must specify at least one file name")
	}
	bx := box.NewFromFlags(c)
	if bx.ConfigRO {
		return fmt.Errorf(roError)
	}
	err := bx.FileRemove(c.Args().Slice())
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdInfo(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	err := bx.Info()
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdInit(c *cli.Context) error {
	if c.Args().Len() > 1 {
		return fmt.Errorf("This command takes one or two arguments")
	}
	bx := box.NewUninitialized(c)
	if bx.ConfigRO {
		return fmt.Errorf(roError)
	}
	err := bx.Init(c.Args().First(), c.String("vcs"))
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdReencrypt(c *cli.Context) error {
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
	err := bx.Reencrypt(c.Args().Slice(),
		c.Bool("overwrite"),
		pauseNeeded,
	)
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdShred(c *cli.Context) error {
	if err := allOrSomeFiles(c); err != nil {
		return err
	}
	bx := box.NewFromFlags(c)
	err := bx.Shred(c.Args().Slice())
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

func cmdStatus(c *cli.Context) error {
	if c.Bool("all") && c.Args().Present() {
		return fmt.Errorf("Can not specify filenames and --all")
	}
	bx := box.NewFromFlags(c)
	err := bx.Status(c.Args().Slice(), c.Bool("name-only"), c.String("type"))
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}

// These are "secret" commands used by the integration tests.

func testingInit(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf("No args required")
	}

	logDebug := bblog.GetDebug(c.Bool("debug"))
	logDebug.Printf(
		"c.String(vcs) reports %q\n",
		c.String("vcs"),
	)
	bx := box.NewForTestingInit(c.String("vcs"))
	if bx.ConfigRO {
		return fmt.Errorf(roError)
	}
	err := bx.TestingInitRepo()
	if err != nil {
		return err
	}
	return bx.Vcs.FlushCommits()
}
