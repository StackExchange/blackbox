package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "2.0.0"
	app.Usage = "Maintain encrypted files in a VCS (Git, Hg, Svn)"

	app.Commands = []cli.Command{
		{
			Name:    "initialize",
			Aliases: []string{"init"},
			Usage:   "Runs blackbox_initialize",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_initialize", c.Args().First())
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e", "ed"},
			Usage:   "Runs blackbox_edit ",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_edit", c.Args().First())
			},
		},
		{
			Name:    "decrypt",
			Aliases: []string{"de"},
			Usage:   "Runs blackbox_edit_start",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_edit_start", c.Args().First())
			},
			// TODO(tlim): Add --all flag to run blackbox_decrypt_all_files
			// TODO(tlim): Add --non-interactive to run blackbox_postdeploy
		},
		{
			Name:    "encrypt",
			Aliases: []string{"en"},
			Usage:   "Runs blackbox_edit_end",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_edit_end", c.Args().First())
			},
		},
		{
			Name:  "reencrypt",
			Usage: "Runs blackbox_update_all_files",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_update_all_files", c.Args().First())
			},
		},
		{
			Name:  "cat",
			Usage: "Runs blackbox_cat",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_cat", c.Args().First())
			},
		},
		{
			Name:  "diff",
			Usage: "Runs blackbox_diff",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_diff", c.Args().First())
			},
		},
		{
			Name:  "shredall",
			Usage: "Runs blackbox_shred_all_files",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_shred_all_files", c.Args().First())
			},
		},
		{
			Name:  "whatsnew",
			Usage: "Runs blackbox_whatsnew",
			Action: func(c *cli.Context) error {
				return RunBash("blackbox_whatsnew", c.Args().First())
			},
		},
		{
			Name:  "admin",
			Usage: "Maintain the list of administrators",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "Runs blackbox_addadmin",
					Action: func(c *cli.Context) error {
						return RunBash("blackbox_addadmin", c.Args().First())
					},
				},
				{
					Name:  "remove",
					Usage: "Runs blackbox_removeadmin",
					Action: func(c *cli.Context) error {
						return RunBash("blackbox_removeadmin", c.Args().First())
					},
				},
				{
					Name:  "list",
					Usage: "Runs blackbox_list_admins",
					Action: func(c *cli.Context) error {
						return RunBash("blackbox_list_admins", c.Args().First())
					},
				},
			},
		},
		{
			Name:    "register",
			Aliases: []string{"reg"},
			Usage:   "Maintain the list of files",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "Runs blackbox_register_new_file",
					Action: func(c *cli.Context) error {
						return RunBash("blackbox_register_new_file", c.Args().First())
					},
				},
				{
					Name:  "remove",
					Usage: "Runs blackbox_deregister_file",
					Action: func(c *cli.Context) error {
						return RunBash("blackbox_deregister_file", c.Args().First())
					},
				},
				{
					Name:  "list",
					Usage: "Runs blackbox_list_admins",
					Action: func(c *cli.Context) error {
						return RunBash("blackbox_list_files", c.Args().First())
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
