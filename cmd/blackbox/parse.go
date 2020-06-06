package main

// All the "blackbox admin" subcommands.

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

// Keep these functions in alphabetical order.

func cmdAdminAdd(c *cli.Context) error {
	if c.Args().Present() {
		return fmt.Errorf(
			"Must specify at least one admin's GnuPG user-id (i.e. email address)")
	}
	bx := box.NewFromFlags(c)
	return bx.AdminAdd(c.Args().Slice())
}

func cmdAdminList(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf("This command takes zero arguments")
	}
	bx := box.NewFromFlags(c)
	return bx.AdminList(bx)
}

func cmdAdminRemove(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdCat(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdDecrypt(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdDiff(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdEdit(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdEncrypt(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdFileAdd(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdFileList(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdFileRemove(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdInfo(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdInit(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdReencrypt(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdShred(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdStatusAll(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdStatusChanged(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
}

func cmdStatusUnchanged(c *cli.Context) error {
	logErr.Println("NOT IMPLEMENTED")
	return nil
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
