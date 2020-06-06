package main

import (
	"fmt"
	"os"
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
