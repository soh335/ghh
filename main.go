package main

import (
	"fmt"
	"github.com/gonuts/commander"
	"os"
	"path/filepath"
)

var ConfigPath = filepath.Join(os.Getenv("HOME"), ".config", os.Args[0])

var mainCmd = &commander.Command{
	UsageLine: os.Args[0],
}

func init() {
	mainCmd.Subcommands = []*commander.Command{
		supportCmd,
		listCmd,
		showCmd,
		createCmd,
		editCmd,
		deleteCmd,
		testCmd,
		configCmd,
	}
}

func main() {
	if err := mainCmd.Dispatch(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
