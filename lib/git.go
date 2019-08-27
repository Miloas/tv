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
	return exec.Command("git", "tag", "-am", commitMessage, tag).Run()
}
