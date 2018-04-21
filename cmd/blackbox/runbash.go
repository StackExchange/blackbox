package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// RunBash runs a Bash command.
func RunBash(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	return errors.Wrapf(err, "run_bash:")
}
