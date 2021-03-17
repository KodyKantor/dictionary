package memdb

import (
	"github.com/hashicorp/go-memdb"
	"github.com/kodykantor/dictionary/pkg/metadb"
)

type MemDB struct {
	db        *memdb.MemDB
	dictTable string
}

func (m *MemDB) InitDB() {
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
}

func (m *MemDB) PutDefinition(def *metadb.Def) error {
	txn := m.db.Txn(true)
	txn.Insert(m.dictTable, def)
	txn.Commit()
	return nil
}

func (m *MemDB) GetDefinition(word string) (*metadb.Def, error) {
	txn := m.db.Txn(false)
	def, err := txn.First(m.dictTable, "id", word)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	if def != nil {
		return def.(*metadb.Def), nil
	}
	return nil, nil
}
