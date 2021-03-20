package metadb_test

import (
	"testing"

	"github.com/kodykantor/dictionary/pkg/metadb"
	"github.com/kodykantor/dictionary/pkg/metadb/memdb"
	"github.com/kodykantor/dictionary/pkg/metadb/memmap"
	"github.com/stretchr/testify/assert"
)

func TestDefinition(t *testing.T) {
	var m metadb.MetaDB
	for _, db := range []string{"memdb", "memmap"} {
		switch db {
		case "memdb":
			m = &memdb.MemDB{}
		case "memmap":
			m = &memmap.MemMap{}
		}

		m.InitDB()

		// Define a word and read it back.
		testDef := &metadb.Def{
			Word:       "milkshake",
			Definition: "a cool summer treat",
		}
		err := m.PutDefinition(testDef)
		assert.Nil(t, err)

		def, err := m.GetDefinition(testDef.Word)
		assert.Nil(t, err)
		assert.Equal(t, testDef, def)

		// Get a word that isn't defined.
		def, err = m.GetDefinition("notdefined")
		assert.Nil(t, err)
		assert.Nil(t, def)

		// Define a word with weird characters and read it back.
		testDef = &metadb.Def{
			Word:       "$$$",
			Definition: "&&&",
		}
		err = m.PutDefinition(testDef)
		assert.Nil(t, err)

		def, err = m.GetDefinition(testDef.Word)
		assert.Nil(t, err)
		assert.Equal(t, testDef, def)
	}

}
