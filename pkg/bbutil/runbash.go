package bbutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunBash runs a Bash command.
func RunBash(command string, args ...string) error {
	//if dryRun {
	//	fmt.Printf("DRY_RUN: Would run exec.Command(%v, %v)\n", command, args)
	//	return nil
	//}
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	return fmt.Errorf("run_bash: %w", err)
}
