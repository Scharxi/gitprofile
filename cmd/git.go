package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func runGitCommand(args ...string) ([]byte, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("git.exe", args...)
	} else {
		cmd = exec.Command("git", args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Error(), strings.TrimSpace(string(output)))
	}
	return output, nil
}
