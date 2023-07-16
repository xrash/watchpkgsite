package git

import (
	"fmt"
	"os/exec"
	"regexp"
)

type BranchSyncState string

const (
	UpToDate BranchSyncState = "up-to-date"
	Behind   BranchSyncState = "behind"
	Ahead    BranchSyncState = "ahead"
)

type GitStatusResult struct {
	LocalBranch  string
	RemoteBranch string
	SyncState    BranchSyncState
	Raw          []byte
}

func GitStatus(runopts *RunOptions) (*GitStatusResult, error) {
	cmd := exec.Command("git", "status")

	if runopts.Dir != "" {
		cmd.Dir = runopts.Dir
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git status: %w", err)
	}

	localBranch, err := parseLocalBranch(output)
	if err != nil {
		return nil, fmt.Errorf("error parsing local branch: %w", err)
	}

	syncState, remoteBranch, err := parseRemoteBranchAndSyncState(output)
	if err != nil {
		return nil, fmt.Errorf("error parsing remote branch: %w", err)
	}

	result := &GitStatusResult{
		LocalBranch:  localBranch,
		RemoteBranch: remoteBranch,
		SyncState:    syncState,
		Raw:          output,
	}

	return result, nil
}

func parseLocalBranch(input []byte) (string, error) {
	regex, err := regexp.Compile("^On branch (.+)")
	if err != nil {
		return "", err
	}

	groups := regex.FindAllSubmatch(input, -1)

	if len(groups) < 1 || len(groups[0]) < 2 {
		return "", fmt.Errorf("couldn't find expected groups, found %v instead", groups)
	}

	return string(groups[0][1]), nil
}

func parseRemoteBranchAndSyncState(input []byte) (BranchSyncState, string, error) {
	upToDateRegex, err := regexp.Compile("Your branch is up to date with '(.+)'")
	if err != nil {
		return "", "", err
	}

	behindRegex, err := regexp.Compile("Your branch is behind '(.+)'")
	if err != nil {
		return "", "", err
	}

	aheadRegex, err := regexp.Compile("Your branch is ahead of '(.+)'")
	if err != nil {
		return "", "", err
	}

	upToDateGroups := upToDateRegex.FindAllSubmatch(input, -1)
	behindGroups := behindRegex.FindAllSubmatch(input, -1)
	aheadGroups := aheadRegex.FindAllSubmatch(input, -1)

	if len(upToDateGroups) < 1 && len(behindGroups) < 1 && len(aheadGroups) < 1 {
		return "", "", fmt.Errorf("couldn't find expected groups, found %v %v %v nstead", upToDateGroups, behindGroups, aheadGroups)
	}

	if len(upToDateGroups) >= 1 && len(upToDateGroups[0]) >= 2 {
		return UpToDate, string(upToDateGroups[0][1]), nil
	}

	if len(behindGroups) >= 1 && len(behindGroups[0]) >= 2 {
		return Behind, string(behindGroups[0][1]), nil
	}

	if len(aheadGroups) >= 1 && len(aheadGroups[0]) >= 2 {
		return Ahead, string(aheadGroups[0][1]), nil
	}

	return "", "", fmt.Errorf("unreachable")
}
