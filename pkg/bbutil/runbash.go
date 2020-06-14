package bbutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	if err != nil {
		return fmt.Errorf("RunBash err=%w", err)
	}
	return nil
}

// RunBash runs a Bash command.
func RunBashOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("RunBashOutput err=%w", err)
	}
	return string(out), err
}
