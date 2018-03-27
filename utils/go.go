package utils

import (
	"encoding/json"
	"fmt"
	"strings"
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

func GoDeps(wkdir string) ([]string, error) {

	cmdArgs := append([]string{"list", "--json"})

	result, err := ExecCommandWD("go", wkdir, cmdArgs...)

	if err != nil {
		err = fmt.Errorf("get go deps failure: %s, output: %s", err.Error(), string(result))
		return nil, err
	}

	goList := struct {
		Imports []string
		Deps    []string
	}{}

	err = json.Unmarshal(result, &goList)
	if err != nil {
		return nil, err
	}

	return append(goList.Imports, goList.Deps...), nil
}

func GoRoot() string {

	result, err := ExecCommand("go", "env", "GOROOT")
	if err != nil {
		return ""
	}

	return strings.TrimSuffix(string(result), "\n")
}
