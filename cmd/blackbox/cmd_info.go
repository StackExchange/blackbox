package main

// All the "blackbox admin" subcommands.

import (
	"fmt"

	"github.com/StackExchange/blackbox/pkg/bbutil"
	"github.com/urfave/cli"
)

func cmdInfo(c *cli.Context) error {

	// GPG version
	// VCS name
	// keys directory

	bbu, err := bbutil.New()
	if err != nil {
		return err
	}

	fmt.Print("VCS:\n")
	fmt.Printf("\tName: %q\n", bbu.Vcs.Name())
	fmt.Printf("\tRepoBaseDir: %q\n", bbu.Vcs.RepoBaseDir())
	fmt.Print("REPO:\n")
	fmt.Printf("\tRepoBaseDir: %q\n", bbu.RepoBaseDir)
	fmt.Printf("\tBlackboxConfigDir: %q\n", bbu.BlackboxConfigDir)

	return nil
}
