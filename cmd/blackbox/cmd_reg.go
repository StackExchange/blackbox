package main

import (
	"fmt"

	"github.com/StackExchange/blackbox/pkg/bbutil"
	"github.com/urfave/cli"
)

func cmdRegList(c *cli.Context) error {
	if len(c.Args()) != 0 {
		fmt.Fprintf(c.App.Writer, "ERROR: Command does not take any arguments\n")
		return nil
	}

	bbu, err := bbutil.New()
	if err != nil {
		return err
	}
	names, err := bbu.RegisteredFiles()
	if err != nil {
		return err
	}
	for _, item := range names {
		fmt.Println(item.Name)
	}
	return nil
}

func cmdRegStatus(c *cli.Context) error {

	if len(c.Args()) != 0 {
		fmt.Fprintf(c.App.Writer, "ERROR: Command does not take any arguments\n")
		return nil
	}

	bbu, err := bbutil.New()
	if err != nil {
		return err
	}
	names, err := bbu.RegisteredFiles()
	if err != nil {
		return err
	}

	for _, item := range names {
		s := bbutil.FileStatus(bbu.RepoBaseDir, item.Name)
		fmt.Printf("%s\t%s\n", s, item.Name)
	}
	return nil
}
