package main

import (
	"github.com/urfave/cli/v2"
)

func flags() *cli.App {
	app := cli.NewApp()
	app.Version = "2.0.0"
	app.Usage = "Maintain encrypted files in a VCS (Git, Hg, Svn)"

	//	app.Flags = []cli.Flag{
	//		&cli.BoolFlag{
	//			Name:        "dry-run",
	//			Aliases:     []string{"n"},
	//			Usage:       "show what would have been done",
	//			Destination: &dryRun,
	//		},
	//	}

	app.Commands = []*cli.Command{

		// List items in the order they appear in the help menu.

		{
			Name:    "decrypt",
			Aliases: []string{"de", "start"},
			Usage:   "Decrypt file(s)",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "all", Usage: "All registered files"},
				&cli.BoolFlag{Name: "bulk", Usage: "Do not prompt to start gpg-agent"},
				&cli.StringFlag{Name: "group", Usage: "Set group ownership"},
				&cli.StringFlag{Name: "overwrite", Usage: "Overwrite plaintext if it exists"},
			},
			Action: func(c *cli.Context) error { return cmdDecrypt(c) },
		},

		{
			Name:    "encrypt",
			Aliases: []string{"en", "end"},
			Usage:   "Encrypts file(s)",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "leave", Usage: "Do not remove plaintext version"},
			},
			Action: func(c *cli.Context) error { return cmdEncrypt(c) },
		},

		{
			Name:    "edit",
			Aliases: []string{"vi"},
			Usage:   "Runs $EDITOR on file(s) (decrypt if needed)",
			Action:  func(c *cli.Context) error { return cmdEdit(c) },
		},

		{
			Name:   "cat",
			Usage:  "Output plaintext to stderr (decrypt if needed)",
			Action: func(c *cli.Context) error { return cmdCat(c) },
		},

		{
			Name:  "diff",
			Usage: "Diffs against encrypted version",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "all", Usage: "all files"},
			},
			Action: func(c *cli.Context) error { return cmdDiff(c) },
		},

		{
			Name:     "init",
			Category: "ADMINISTRATIVE",
			Usage:    "Initialized blackbox for this repo",
			Action:   func(c *cli.Context) error { return cmdInit(c) },
		},

		{
			Name:     "admin",
			Category: "ADMINISTRATIVE",
			Usage:    "Add/list/remove administrators",
			Subcommands: []*cli.Command{
				{
					Name:   "add",
					Usage:  "Adds admin(s)",
					Action: func(c *cli.Context) error { return cmdAdminAdd(c) },
				},
				{
					Name:   "list",
					Usage:  "Lists admins",
					Action: func(c *cli.Context) error { return cmdAdminList(c) },
				},
				{
					Name:   "remove",
					Usage:  "Remove admin(s)",
					Action: func(c *cli.Context) error { return cmdAdminRemove(c) },
				},
			},
		},

		{
			Name:     "file",
			Category: "ADMINISTRATIVE",
			Usage:    "Add/list/remove files from the registry",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "Registers file with the system",
					Flags: []cli.Flag{
						&cli.BoolFlag{Name: "leave", Usage: "Do not remove plaintext version"},
					},
					Action: func(c *cli.Context) error { return cmdFileAdd(c) },
				},
				{
					Name:   "list",
					Usage:  "Lists the registered files",
					Action: func(c *cli.Context) error { return cmdFileList(c) },
				},
				{
					Name:   "remove",
					Usage:  "Deregister file from the system",
					Action: func(c *cli.Context) error { return cmdFileRemove(c) },
				},
			},
		},

		{
			Name:     "info",
			Category: "DEBUG",
			Usage:    "Report what we know about this repo",
			Action:   func(c *cli.Context) error { return cmdInfo(c) },
		},

		{
			Name:  "shred",
			Usage: "Shred the plaintext",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "all", Usage: "All registered files"},
			},
			Action: func(c *cli.Context) error { return cmdShred(c) },
		},

		{
			Name:     "status",
			Category: "ADMINISTRATIVE",
			Usage:    "Print status of files",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "name-only", Usage: "Show only names of the files"},
				&cli.BoolFlag{Name: "all", Usage: "All registered files"},
				&cli.StringFlag{Name: "type", Usage: "only list if status matching this string"},
			},
			Action: func(c *cli.Context) error { return cmdStatus(c) },
		},

		{
			Name:     "reencrypt",
			Usage:    "Decrypt then re-encrypt files",
			Category: "ADMINISTRATIVE",
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "all", Usage: "All registered files"},
			},
			Action: func(c *cli.Context) error { return cmdReencrypt(c) },
		},

		//

	}

	return app
}
