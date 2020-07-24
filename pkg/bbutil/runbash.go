package bbutil

import (
	"bytes"
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
		return fmt.Errorf("RunBash cmd=%q err=%w", command, err)
	}
	return nil
}

// RunBashOutput runs a Bash command, captures output.
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

// RunBashOutputSilent runs a Bash command, captures output, discards stderr.
func RunBashOutputSilent(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	// Leave cmd.Stderr unmodified and stderr is discarded.
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("RunBashOutputSilent err=%w", err)
	}
	return string(out), err
}

// RunBashInput runs a Bash command, sends input on stdin.
func RunBashInput(input string, command string, args ...string) error {

	cmd := exec.Command(command, args...)
	cmd.Stdin = bytes.NewBuffer([]byte(input))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("RunBashInput err=%w", err)
	}
	return nil
}

// RunBashInputOutput runs a Bash command, sends input on stdin.
func RunBashInputOutput(input []byte, command string, args ...string) ([]byte, error) {

	cmd := exec.Command(command, args...)
	cmd.Stdin = bytes.NewBuffer(input)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("RunBashInputOutput err=%w", err)
	}
	return out, nil
}
