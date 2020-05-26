package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, file := range files {
		if !file.Mode().IsRegular() {
			continue
		}

		if strings.ContainsRune(file.Name(), '=') {
			continue
		}

		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			return env, err
		}

		// Prepare env value
		value := strings.Split(string(content), "\n")[0]
		value = strings.ReplaceAll(value, string(0), "\n")
		for _, c := range [...]string{"\t", " "} {
			value = strings.TrimRight(value, c)
		}

		env[file.Name()] = value
	}

	return env, nil
}
