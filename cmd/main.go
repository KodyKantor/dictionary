package main

import (
	"log"
	"net/http"

	"github.com/kodykantor/dictionary/pkg/dictionary"
)

func main() {
	dict := dictionary.Dictionary{}
	dict.Open("memdb") // memdb for development.

	http.HandleFunc("/definition", dict.HandleDefinition)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
