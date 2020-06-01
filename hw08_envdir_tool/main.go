package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("programm syntax: <env directory> <external programm> [<external programm args...>]")
	}

	envDir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}

	for k := range env {
		if _, ok := os.LookupEnv(k); !ok {
			continue
		}
		if err := os.Unsetenv(k); err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(RunCmd(cmd, env))
}
