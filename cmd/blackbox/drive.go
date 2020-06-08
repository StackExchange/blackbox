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
	if c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	return bx.Init()
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

// func cmdRegList(c *cli.Context) error {
// 	if len(c.Args()) != 0 {
// 		fmt.Fprintf(c.App.Writer, "ERROR: Command does not take any arguments\n")
// 		return nil
// 	}
//
// 	bbu, err := bbutil.New()
// 	if err != nil {
// 		return err
// 	}
// 	names, err := bbu.RegisteredFiles()
// 	if err != nil {
// 		return err
// 	}
// 	for _, item := range names {
// 		fmt.Println(item.Name)
// 	}
// 	return nil
// }

// func cmdRegStatus(c *cli.Context) error {
//
// 	if len(c.Args()) != 0 {
// 		fmt.Fprintf(c.App.Writer, "ERROR: Command does not take any arguments\n")
// 		return nil
// 	}
//
// 	bbu, err := bbutil.New()
// 	if err != nil {
// 		return err
// 	}
// 	names, err := bbu.RegisteredFiles()
// 	if err != nil {
// 		return err
// 	}
//
// 	for _, item := range names {
// 		s := bbutil.FileStatus(bbu.RepoBaseDir, item.Name)
// 		fmt.Printf("%s\t%s\n", s, item.Name)
// 	}
// 	return nil
// }
