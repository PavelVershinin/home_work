package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("can't read dir: %w", err)
	}

	env := make(Environment)

	for _, file := range files {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}
		val, err := content(dir, file)
		if err != nil {
			return nil, err
		}
		env[file.Name()] = val
	}

	return env, nil
}

func content(dir string, file os.FileInfo) (string, error) {
	if file.Size() == 0 {
		return "", nil
	}
	b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
	if err != nil {
		return "", fmt.Errorf("can't read file: %w", err)
	}
	b = bytes.Split(b, []byte("\n"))[0]
	b = bytes.ReplaceAll(b, []byte("\x00"), []byte("\n"))
	b = bytes.TrimRight(b, `\s`)

	return string(b), nil
}
