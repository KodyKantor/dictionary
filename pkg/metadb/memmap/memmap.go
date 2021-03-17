package memmap

import "github.com/kodykantor/dictionary/pkg/metadb"

type MemMap struct {
	dict map[string]string
}

func (m *MemMap) InitDB() {
	m.dict = make(map[string]string)
}

func (m *MemMap) PutDefinition(def *metadb.Def) error {
	m.dict[def.Word] = def.Definition
	return nil
}

func (m *MemMap) GetDefinition(word string) (*metadb.Def, error) {
	def := m.dict[word]
	if def == "" {
		return nil, nil
	}
	return &metadb.Def{
		Word:       word,
		Definition: def,
	}, nil
}
