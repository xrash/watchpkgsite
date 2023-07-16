package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderbookGetAmountOut(t *testing.T) {

	{
		input := `On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean
`

		localBranch, err := parseLocalBranch([]byte(input))

		assert.Equal(t, nil, err, "err should be nil")
		assert.Equal(t, "main", localBranch, "localBranch should be main")
	}

	{
		input := `On branch main
Your branch is behind 'origin/main' by 1 commit, and can be fast-forwarded.
  (use "git pull" to update your local branch)

nothing to commit, working tree clean
`

		syncState, remoteBranch, err := parseRemoteBranchAndSyncState([]byte(input))

		assert.Equal(t, nil, err, "err should be nil")
		assert.Equal(t, "origin/main", remoteBranch, "remoteBranch should be origin/main")
		assert.Equal(t, Behind, syncState, "syncState should be Behind")
	}

	{
		input := `On branch main
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   app/root/cmd.go

no changes added to commit (use "git add" and/or "git commit -a")
`

		syncState, remoteBranch, err := parseRemoteBranchAndSyncState([]byte(input))

		assert.Equal(t, nil, err, "err should be nil")
		assert.Equal(t, "origin/main", remoteBranch, "remoteBranch should be origin/main")
		assert.Equal(t, Ahead, syncState, "syncState should be Behind")
	}

}
