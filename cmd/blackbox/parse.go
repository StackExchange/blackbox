package main

// Processes the flags and arguments and calls the appropriate
// business logic.

import (
	"fmt"
	"log"
	"os"

	"github.com/StackExchange/blackbox/pkg/box"
	"github.com/urfave/cli/v2"
)

var logErr *log.Logger

func init() {
	if logErr == nil {
		logErr = log.New(os.Stderr, "", 0)
	}
}

func allOrSomeFiles(c *cli.Context) error {
	if c.Bool("all") {
		if c.Args().Present() {
			return fmt.Errorf("Can not specify filenames and --all")
		}
	} else {
		if !c.Args().Present() {
			return fmt.Errorf("Must specify at least one file name")
		}
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
	bulk := false
	if c.Bool("all") {
		bulk = c.Bool("bulk") // Only applies to --all
	}
	bx := box.NewFromFlags(c)
	return bx.Decrypt(c.Args().Slice(),
		c.Bool("overwrite"),
		bulk,
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
	if err := allOrSomeFiles(c); err != nil {
		return err
	}

	mode := box.Itemized

	if c.Bool("all") {
		if c.Args().Present() {
			return fmt.Errorf("Can not specify filenames and --all")
		}
		if mode != box.Itemized {
			return fmt.Errorf("--all can not be mixed with other flags")
		}
		mode = box.All
	}

	if c.Bool("changed") {
		if c.Args().Present() {
			return fmt.Errorf("Can not specify filenames and --changed")
		}
		if mode != box.Itemized {
			return fmt.Errorf("--changed can not be mixed with other flags")
		}
		mode = box.Changed
	}

	if c.Bool("unchanged") {
		if c.Args().Present() {
			return fmt.Errorf("Can not specify filenames and --unchanged")
		}
		if mode != box.Itemized {
			return fmt.Errorf("--unchanged can not be mixed with other flags")
		}
		mode = box.Unchanged
	}

	bx := box.NewFromFlags(c)
	return bx.Status(c.Args().Slice(), mode, c.Bool("name-only"))
}

//func cmdInfo(c *cli.Context) error {
//
//	// GPG version
//	// VCS name
//	// keys directory
//
//	bbu, err := bbutil.New()
//	if err != nil {
//		return err
//	}
//
//	fmt.Print("VCS:\n")
//	fmt.Printf("\tName: %q\n", bbu.Vcs.Name())
//	fmt.Printf("\tRepoBaseDir: %q\n", bbu.Vcs.RepoBaseDir())
//	fmt.Print("REPO:\n")
//	fmt.Printf("\tRepoBaseDir: %q\n", bbu.RepoBaseDir)
//	fmt.Printf("\tBlackboxConfigDir: %q\n", bbu.BlackboxConfigDir)
//
//	return nil
//}

//func cmdAdminList(c *cli.Context) error {
//
//	if len(c.Args()) != 0 {
//		fmt.Fprintln(c.App.Writer, "ERROR: 'blackbox admin list' does not take any arguments")
//		return nil
//	}
//
//	bbu, err := bbutil.New()
//	if err != nil {
//		return err
//	}
//	names, err := bbu.Administrators()
//	if err != nil {
//		return err
//	}
//	for _, item := range names {
//		fmt.Println(item.Name)
//	}
//	return nil
//}

// func cmdDecrypt(allFiles bool, filenames []string, group string) error {
// 	bbu, err := bbutil.New()
// 	if err != nil {
// 		return err
// 	}
//
// 	// prepare_keychain
//
// 	fnames, valid, err := bbu.FileIterator(allFiles, filenames)
// 	if err != nil {
// 		return errors.Wrap(err, "decrypt")
// 	}
// 	for i, filename := range fnames {
// 		if valid[i] {
// 			bbu.DecryptFile(filename, group, true)
// 		} else {
// 			fmt.Fprintf(os.Stderr, "SKIPPING: %q\n", filename)
// 		}
// 	}
//
// 	return nil
// }

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
