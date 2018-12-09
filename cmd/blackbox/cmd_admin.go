package main

// All the "blackbox admin" subcommands.

import (
	"fmt"

	"github.com/StackExchange/blackbox/pkg/bbutil"
	"github.com/urfave/cli"
)

func cmdAdminList(c *cli.Context) error {

	if len(c.Args()) != 0 {
		fmt.Fprintln(c.App.Writer, "ERROR: 'blackbox admin list' does not take any arguments")
		return nil
	}

	bbu, err := bbutil.New()
	if err != nil {
		return err
	}
	names, err := bbu.Administrators()
	if err != nil {
		return err
	}
	for _, item := range names {
		fmt.Println(item.Name)
	}
	return nil
}
