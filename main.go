package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-memdb"
)

// XXX db operations should all be able to return errors.
type metaDB interface {
	InitDB()                                 // Initialize the database.
	PutDefinition(def *Def)                  // Insert the given definition into the database.
	GetDefinition(word string) (*Def, error) // Get the given word's definition from the database.
}

type memDB struct {
	db        *memdb.MemDB
	dictTable string
}

// Dictionary stores definitions of words. These definitions can be overwritten or retrieved.
type Dictionary struct {
	db metaDB
}

func (m *memDB) InitDB() {
	m.dictTable = "dictionary"

	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			m.dictTable: &memdb.TableSchema{
				Name: m.dictTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Word"},
					},
				},
			},
		},
	}

	// Create a new database
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	m.db = db

	m.PutDefinition(&Def{
		Word:       "prometheus",
		Definition: "the beginning",
	})
}

func (m *memDB) PutDefinition(def *Def) {
	txn := m.db.Txn(true)
	txn.Insert(m.dictTable, def)
	txn.Commit()
}

func (m *memDB) GetDefinition(word string) (*Def, error) {
	txn := m.db.Txn(false)
	def, err := txn.First(m.dictTable, "id", word)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	if def != nil {
		return def.(*Def), nil
	}
	return nil, nil
}

func (d *Dictionary) Open(backend string) {
	// Set up the database if it's not already configured.
	// XXX implement a durable metadata backend.
	if d.db == nil {
		switch backend {
		case "memdb":
			m := memDB{}
			d.db = &m
		default:
			panic("unknown db type")
		}

		d.db.InitDB()
	}
}

// Def represents a single word and its definition.
type Def struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

func (d *Dictionary) getDefinition(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")
	if word == "" {
		http.Error(w, "a word must be requested", http.StatusBadRequest)
		return
	}

	myDef, err := d.db.GetDefinition(word)
	if err != nil {
		http.Error(w, "server GetDefinition error", http.StatusInternalServerError)
		return
	}

	if myDef == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s is not in the dictionary\n", word)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s: %s\n", myDef.Word, myDef.Definition)
	}
}

func (d *Dictionary) putDefinition(w http.ResponseWriter, r *http.Request) {
	// XXX set a buffer size limit.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "server ioutil error", http.StatusInternalServerError)
		return
	}

	myDef := Def{}

	err = json.Unmarshal(body, &myDef)
	if err != nil {
		http.Error(w, "bad json payload", http.StatusBadRequest)
		return
	}

	d.db.PutDefinition(&myDef)
	w.WriteHeader(http.StatusCreated)
}

// HandleDefinition does post-routing of requests depending on the request type.
func (d *Dictionary) HandleDefinition(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		d.getDefinition(w, r)
	case http.MethodPut:
		d.putDefinition(w, r)
	default:
		http.Error(w, "unsupported operation", http.StatusBadRequest)
	}
	return
}

func main() {
	dict := Dictionary{}
	dict.Open("memdb") // memdb for development.

	http.HandleFunc("/definition", dict.HandleDefinition)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
