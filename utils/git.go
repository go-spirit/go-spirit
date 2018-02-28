package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func GitRepoName(gitUrl string) string {
	found := regexp.MustCompile(`([^/]+)\.git$`).FindAllString(gitUrl, 1)

	if len(found) == 0 {
		return ""
	}

	return strings.TrimSuffix(found[0], ".git")
}

func GitCheckout(repoDir, revision string) (err error) {
	checkoutArgs := []string{"-C", repoDir, "checkout", revision}

	var output []byte
	output, err = ExecCommand("git", checkoutArgs...)
	if err != nil {
		err = fmt.Errorf("checkout revision failure, dir: %s, error: %s\n%s\n", repoDir, output, err)
		return
	}

	return
}

func DirExist(dir string) (bool, error) {
	fi, errF := os.Stat(dir)
	if errF != nil {
		if os.IsNotExist(errF) {
			return false, nil
		}
		return false, errF
	} else if !fi.IsDir() {
		return false, fmt.Errorf("%s is not a dir", dir)
	}

	return true, nil
}

func GitDetached(wkdir string) (bool, error) {
	result, err := ExecCommand("git", "-C", wkdir, "status", "-b")

	if err != nil {
		err = fmt.Errorf("get git status failure: %s", err.Error())
		return false, err
	}

	return strings.Index(string(result), "detached") > 0, nil
}

func GitPull(wkdir string, args ...string) error {

	cmdArgs := append([]string{"-C", wkdir, "pull"}, args...)

	output, err := ExecCommand("git", cmdArgs...)

	if err != nil {
		err = fmt.Errorf("git pull failure (%s): %s\n%s", wkdir, string(output), err.Error())
		return err
	}

	return nil
}

func GitClone(wkdir, url string, args ...string) error {

	cmdArgs := append([]string{"-C", wkdir, "clone"}, args...)
	cmdArgs = append(cmdArgs, url)

	output, err := ExecCommand("git", cmdArgs...)

	if err != nil {
		err = fmt.Errorf("git clone failure: %s\n%s", string(output), err.Error())
		return err
	}

	return nil
}
