package utils

import (
	"fmt"
)

func GoGet(url string, args ...string) error {

	cmdArgs := append([]string{"get"}, args...)

	cmdArgs = append(cmdArgs, url)

	result, err := ExecCommand("go", cmdArgs...)

	if err != nil {
		err = fmt.Errorf("go get failure: %s, output: %s", err.Error(), string(result))
		return err
	}

	return nil
}
