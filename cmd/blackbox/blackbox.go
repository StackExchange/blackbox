package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var dryRun bool

func main() {
	app := cli.NewApp()
	app.Version = "2.0.0"
	app.Usage = "Maintain encrypted files in a VCS (Git, Hg, Svn)"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "dry-run, n",
			Usage:       "show what would have been done",
			Destination: &dryRun,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "info",
			Category: "DEBUG",
			Usage:    "Report what we know about this repo",
			Action:   func(c *cli.Context) error { return cmdInfo(c) },
		},

		{
			Name:     "initialize",
			Aliases:  []string{"init"},
			Category: "GETTING STARTED",
			Usage:    "Runs blackbox_initialize",
			Action:   func(c *cli.Context) error { return RunBash("blackbox_initialize", c.Args().First()) },
		},
		{
			Name:    "edit",
			Aliases: []string{"e", "ed"},
			Usage:   "Runs blackbox_edit ",
			Action:  func(c *cli.Context) error { return RunBash("blackbox_edit", c.Args().First()) },
		},
		{
			Name:    "decrypt",
			Aliases: []string{"de", "start"},
			Usage:   "Runs blackbox_edit_start",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "all, a", Usage: "Runs blackbox_decrypt_all_files"},
				cli.BoolFlag{Name: "non-interactive", Usage: "Runs blackbox_postdeploy"},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("all") {
					if c.Bool("non-interactive") {
						return RunBash("blackbox_postdeploy", c.Args().First())
					}
					return RunBash("blackbox_decrypt_all_files", c.Args().First())
				}
				return RunBash("blackbox_edit_start", c.Args().First())
			},
		},
		{
			Name:    "ndecrypt",
			Aliases: []string{"de", "start"},
			Usage:   "Runs blackbox_edit_start",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "All registered files",
				},
				cli.BoolFlag{
					Name:  "non-interactive, b",
					Usage: "Do not set up gpg-agent",
				},
				cli.StringFlag{
					Name:  "set-group, g",
					Usage: "Set group ownership",
				},
			},
			Action: func(c *cli.Context) error {
				if !c.Bool("non-interactive") {
					// gpg_agent_notice
				}
				if c.Bool("all") {
					if len(c.Args()) != 0 {
						return errors.New("Can't combined --all and filenames")
					}
				} else {
					if len(c.Args()) == 0 {
						return errors.New("At least one filename required (or --all)")
					}
				}
				return cmdDecrypt(c.Bool("all"), c.Args(), c.String("set-group"))
			},
		},
		{
			Name:    "encrypt",
			Aliases: []string{"en", "end"},
			Usage:   "Runs blackbox_edit_end",
			Action:  func(c *cli.Context) error { return RunBash("blackbox_edit_end", c.Args().First()) },
		},
		{
			Name:   "reencrypt",
			Usage:  "Runs blackbox_update_all_files",
			Action: func(c *cli.Context) error { return RunBash("blackbox_update_all_files", c.Args().First()) },
		},
		{
			Name:   "cat",
			Usage:  "Runs blackbox_cat",
			Action: func(c *cli.Context) error { return RunBash("blackbox_cat", c.Args().First()) },
		},
		{
			Name:   "diff",
			Usage:  "Runs blackbox_diff",
			Action: func(c *cli.Context) error { return RunBash("blackbox_diff", c.Args().First()) },
		},
		{
			Name:   "shredall",
			Usage:  "Runs blackbox_shred_all_files",
			Action: func(c *cli.Context) error { return RunBash("blackbox_shred_all_files", c.Args().First()) },
		},
		{
			Name:   "whatsnew",
			Usage:  "Runs blackbox_whatsnew",
			Action: func(c *cli.Context) error { return RunBash("blackbox_whatsnew", c.Args().First()) },
		},
		{
			Name:     "admin",
			Category: "ADMINISTRATIVE",
			Usage:    "Maintain the list of administrators",
			Subcommands: []cli.Command{
				{
					Name:    "nadd",
					Aliases: []string{"add"},
					Usage:   "Adds adminstrator(s)",
					Action:  func(c *cli.Context) error { return RunBash("blackbox_addadmin", c.Args().First()) },
				},
				{
					Name:   "oadd",
					Usage:  "Runs blackbox_addadmin",
					Action: func(c *cli.Context) error { return RunBash("blackbox_addadmin", c.Args().First()) },
				},
				{
					Name:   "remove",
					Usage:  "Runs blackbox_removeadmin",
					Action: func(c *cli.Context) error { return RunBash("blackbox_removeadmin", c.Args().First()) },
				},
				{
					Name:   "olist",
					Usage:  "Runs blackbox_list_admins",
					Action: func(c *cli.Context) error { return RunBash("blackbox_list_admins", c.Args().First()) },
				},
				{
					Name:    "nlist",
					Aliases: []string{"list"},
					Usage:   "Lists blackbox admins",
					Action:  func(c *cli.Context) error { return cmdAdminList(c) },
				},
			},
		},
		{
			Name:     "file",
			Aliases:  []string{"f"},
			Category: "ADMINISTRATIVE",
			Usage:    "Maintain the list of files",
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "Runs blackbox_register_new_file",
					Action: func(c *cli.Context) error { return RunBash("blackbox_register_new_file", c.Args().First()) },
				},
				{
					Name:   "remove",
					Usage:  "Runs blackbox_deregister_file",
					Action: func(c *cli.Context) error { return RunBash("blackbox_deregister_file", c.Args().First()) },
				},
				{
					Name:   "olist",
					Usage:  "Runs blackbox_list_admins",
					Action: func(c *cli.Context) error { return RunBash("blackbox_list_files", c.Args().First()) },
				},
				{
					Name:    "nlist",
					Aliases: []string{"list"},
					Usage:   "Lists the registered files",
					Action:  func(c *cli.Context) error { return cmdRegList(c) },
				},
				{
					Name:   "status",
					Usage:  "Prints info about registered files",
					Action: func(c *cli.Context) error { return cmdRegStatus(c) },
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
