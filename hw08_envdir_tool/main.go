package main

import (
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	args := os.Args[1:]

	//nolint:gomnd
	if len(args) < 2 {
		log.Fatalln("you must pass at least two arguments!")
	}

	envDir := args[0]
	cmd := args[1:]

	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(RunCmd(cmd, env))
}
