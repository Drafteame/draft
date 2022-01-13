package utils

import (
	"fmt"
	"os/exec"
)

// Runs a system call and returns an error if anything goes wrong,
// including non-zero return status.
//
// argv[0] is the path to the program.
// The arguments are in argv[1:].
func Run(argv ...string) error {
	cmd := exec.Command(argv[0], argv[1:]...)
	stdout, err := cmd.Output()

	if err != nil {
		return err
	}

	// Print the output
	fmt.Println(string(stdout))
	return nil
}
