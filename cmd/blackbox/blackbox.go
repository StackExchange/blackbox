package main

import (
	"fmt"
	"os"

	_ "github.com/StackExchange/blackbox/crypters"
	_ "github.com/StackExchange/blackbox/crypters/_all"
	_ "github.com/StackExchange/blackbox/vcs"
	_ "github.com/StackExchange/blackbox/vcs/_all"
)

var dryRun bool

func main() {
	app := flags()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
