package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

//nolint:godot
// RunCmd runs a command + arguments (cmd) with environment variables from env
func RunCmd(cmd []string, env Environment) (returnCode int) {
	return runCmd(cmd[0], cmd[1:], env, os.Stdin, os.Stdout, os.Stderr)
}

func runCmd(cmd string, args []string, env Environment, stdIn io.Reader, stdOut, stdErr io.Writer) int {
	run := exec.Command(cmd, args...)
	run.Stdin = stdIn
	run.Stdout = stdOut
	run.Stderr = stdErr

	for _, e := range os.Environ() {
		a := strings.Split(e, "=")
		if _, ok := env[a[0]]; ok {
			continue
		}
		run.Env = append(run.Env, e)
	}

	for name, value := range env {
		if len(value) > 0 {
			run.Env = append(run.Env, name+"="+value)
		}
	}

	if err := run.Run(); err != nil {
		log.Println(err)
	}

	return run.ProcessState.ExitCode()
}
