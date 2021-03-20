package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kodykantor/dictionary/pkg/dictionary"
)

var (
	definitionPath = "/definition"
)

// DictionaryServer provides an HTTP abstraction over a Dictionary.
type DictionaryServer struct {
	port   string // tcp port for server to listen on.
	dbType string // type of database for definition storage.
}

func (ds *DictionaryServer) Help() string {
	helpText := `
	Start a dictionary server on port 8080.

	The dictionary server is not daemonized.
	`

	return helpText
}

func (ds *DictionaryServer) Synopsis() string {
	return "Start a dictionary server."
}

func (ds *DictionaryServer) Run(args []string) int {
	dict := dictionary.Dictionary{}
	dict.Open(ds.dbType)

	http.HandleFunc(definitionPath, dict.HandleDefinition)

	fmt.Printf("starting server on port %d...\n", ds.port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", ds.port), nil)
	if err != nil {
		log.Println("error running server:", err)
	}

	return 0
}
