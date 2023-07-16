package git

import (
	"fmt"
	"os/exec"
)

type GitFetchResult struct {
	Raw []byte
}

func GitFetch(runopts *RunOptions) (*GitFetchResult, error) {
	cmd := exec.Command("git", "fetch")

	if runopts.Dir != "" {
		cmd.Dir = runopts.Dir
	}

	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git fetch: %w", err)
	}

	return nil, nil
}
