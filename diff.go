package gitlib

import (
	"os"
	"os/exec"
	"strings"
)

// RunDiff runs `git diff` command and show it.
func RunDiff(path string, isStaged bool) error {
	var cmd *exec.Cmd
	if isStaged {
		cmd = exec.Command("git", "diff", "--staged", path)
	} else {
		cmd = exec.Command("git", "diff", path)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// DiffFile is file that has diff.
type DiffFile struct {
	Status   string
	Path     string
	IsStaged bool
}

// GetDiffFiles gets diff files.
// If isStaged is true, you can get diff files from staging.
// and if false, get diff files from unstaging.
func GetDiffFiles(isStaged bool) ([]DiffFile, error) {
	nameOnly, err := GetDiffNameOnly(isStaged)
	if err != nil {
		return nil, err
	}
	nameStatus, err := GetDiffNameStatus(isStaged)
	if err != nil {
		return nil, err
	}

	// len(nameStatuses) equal len(nameOnlies)
	rowLen := len(nameOnly)

	diffFile := make([]DiffFile, rowLen)
	for i := 0; i < rowLen; i++ {
		path := nameOnly[i]
		status := strings.Fields(nameStatus[i])

		diffFile[i] = DiffFile{
			Status:   status[0],
			Path:     path,
			IsStaged: isStaged,
		}
	}

	return diffFile, nil
}

// GetDiffNameOnly gets result that run `git diff --name-only` git commnad.
// If isStaged is true, add `--staged`.
func GetDiffNameOnly(isStaged bool) ([]string, error) {
	var cmd *exec.Cmd
	if isStaged {
		cmd = exec.Command("git", "diff", "--staged", "--name-only")
	} else {
		exec.Command("git", "add", "-A", "-N").Run()
		cmd = exec.Command("git", "diff", "--name-only")
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(out), "\n")
	// Remove the latest empty row.
	return rows[0 : len(rows)-1], nil
}

// GetDiffNameStatus gets result that run `git diff --name-status` git command.
// If isStaged is true, add `--staged`.
func GetDiffNameStatus(isStaged bool) ([]string, error) {
	var cmd *exec.Cmd
	if isStaged {
		cmd = exec.Command("git", "diff", "--staged", "--name-status")
	} else {
		exec.Command("git", "add", "-A", "-N").Run()
		cmd = exec.Command("git", "diff", "--name-status")
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	rows := strings.Split(string(out), "\n")
	// Remove the latest empty row.
	return rows[0 : len(rows)-1], nil
}
