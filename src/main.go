package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const (
	debugMode bool = true
)

func main() {
	commandMap := map[string]executer{
		"version": versionCommand(),
		"init":    initCommand(),
		"list":    listCommand(),
		"install": installCommand(),
	}

	commandArg := os.Args[1]
	command, ok := commandMap[commandArg]
	if ok {
		err := command.Execute()
		if err != nil {
			if debugMode {
				panic(err)
			} else {
				fmt.Println("Oops an error occured, submit an issue and provide any information")
			}
		}
	} else {
		fmt.Println("Command", commandArg, "not recognized")
	}
}

func executableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	linkPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}

	return path.Dir(linkPath), nil
}
