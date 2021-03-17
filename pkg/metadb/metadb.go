package metadb

// XXX db operations should all be able to return errors.
type MetaDB interface {
	InitDB()                                 // Initialize the database.
	PutDefinition(def *Def)                  // Insert the given definition into the database.
	GetDefinition(word string) (*Def, error) // Get the given word's definition from the database.
}

// Def represents a single word and its definition.
type Def struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}
