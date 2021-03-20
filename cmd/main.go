package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

const (
	appVersion = "1.0.0"
	appName    = "dictionary"

	defaultPort = "8080"
	defaultDB   = "memmap"
)

func serverFactory() (cli.Command, error) {
	return &DictionaryServer{
		port:   defaultPort,
		dbType: defaultDB,
	}, nil
}

func clientFactory() (cli.Command, error) {
	return &DictionaryClient{
		remotePort: defaultPort,
	}, nil
}

func main() {

	c := cli.NewCLI(appName, appVersion)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"server": serverFactory,
		"client": clientFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println("error running command:", err)
	}

	os.Exit(exitStatus)
}
