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

func commitChanges(version string) (string, error) {
	commitMessage, _ := strconv.Unquote("\"" + version + "\"")
	err := exec.Command("git", "add", SemverFileName).Run()
	if err != nil {
		return "", err
	}
	err = exec.Command("git", "commit", "-m", commitMessage, "-n").Run()
	if err != nil {
		return "", err
	}
	return commitMessage, nil
}

func TagVersion(tag string) error {
	commitMessage, err := commitChanges(tag)
	if err != nil {
		return err
	}
	return exec.Command("git", "tag", "-am", commitMessage, tag).Run()
}

func TagVersions(version string, tags []string) error {
	commitMessage, err := commitChanges(version)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		err := exec.Command("git", "tag", "-am", commitMessage, tag).Run()
		if err != nil {
			return err
		}
	}
	return nil
}
