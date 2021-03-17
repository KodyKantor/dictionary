package dictionary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kodykantor/dictionary/pkg/metadb"
	"github.com/kodykantor/dictionary/pkg/metadb/memdb"
	"github.com/kodykantor/dictionary/pkg/metadb/memmap"
)

// Dictionary stores definitions of words. These definitions can be overwritten or retrieved.
type Dictionary struct {
	db metadb.MetaDB
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

	myDef := metadb.Def{}

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

func (d *Dictionary) Open(backend string) {
	// Set up the database if it's not already configured.
	// XXX implement a durable metadata backend.
	if d.db == nil {
		switch backend {
		case "memdb":
			m := memdb.MemDB{}
			d.db = &m
		case "memmap":
			m := memmap.MemMap{}
			d.db = &m
		default:
			panic("unknown db type")
		}

		d.db.InitDB()

		d.db.PutDefinition(&metadb.Def{
			Word:       "prometheus",
			Definition: "the beginning",
		})
	}
}
