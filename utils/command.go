package utils

import (
	"os"
	"os/exec"
)

func ExecCommand(name string, args ...string) (result []byte, err error) {

	result, err = exec.Command(name, args...).CombinedOutput()

	return
}

func ExecCommandSTD(name string, args ...string) (err error) {

	cmd := exec.Command(name, args...)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Start()
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	return
}

func ExecCommandWD(name string, dir string, args ...string) (err error) {

	cmd := exec.Command(name, args...)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Dir = dir

	err = cmd.Start()
	if err != nil {
		return
	}

	err = cmd.Wait()
	if err != nil {
		return
	}

	return
}
