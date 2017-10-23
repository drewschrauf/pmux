package util

import (
	"os/exec"
	"regexp"
	"strings"
)

func GitBranch(dir string) (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "--quiet", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	} else {
		ref := strings.TrimSpace(string(out))
		re := regexp.MustCompile("^refs/heads/(.*)$")
		return re.FindStringSubmatch(ref)[1], nil
	}
}

func GitAhead(dir string) (bool, error) {
	branch, err := GitBranch(dir)
	if err != nil {
		return false, err
	}
	cmd := exec.Command("git", "rev-list", branch+"@{upstream}..HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	revs := strings.TrimSpace(string(out))
	return revs != "", nil
}

func GitBehind(dir string) (bool, error) {
	branch, err := GitBranch(dir)
	if err != nil {
		return false, err
	}

	cmd := exec.Command("git", "remote", "update")
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		return false, err
	}

	cmd = exec.Command("git", "rev-list", "HEAD.."+branch+"@{upstream}")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	revs := strings.TrimSpace(string(out))
	return revs != "", nil
}

func GitDirty(dir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	diff := strings.TrimSpace(string(out))
	return diff != "", nil
}
