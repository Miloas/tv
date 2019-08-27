package tv

import (
	"os/exec"
	"strconv"
)

func IsClean() bool {
	out, err := exec.Command("git", "status", "--untracked", "--short").Output()
	if err != nil {
		panic(err)
	}
	return string(out) == ""
}

func TagVersion(tag string) error {
	commitMessage, _ := strconv.Unquote("\"" + tag + "\"")
	err := exec.Command("git", "tag", "-am", commitMessage, tag).Run()
	if err != nil {
		return err
	}
	err = exec.Command("git", "add", "semver.json").Run()
	if err != nil {
		return err
	}
	return exec.Command("git", "commit", "-m", commitMessage).Run()
}
