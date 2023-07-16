package git

import (
	"fmt"
	"os/exec"
)

type GitMergeResult struct {
	Raw []byte
}

func GitMerge(runopts *RunOptions, _args []string) (*GitMergeResult, error) {
	args := append([]string{"merge"}, _args...)
	cmd := exec.Command("git", args...)

	if runopts.Dir != "" {
		cmd.Dir = runopts.Dir
	}

	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git merge: %w", err)
	}

	return nil, nil
}
