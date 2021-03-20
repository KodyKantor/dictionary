package dictionary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	d := &Dictionary{}

	err := d.Open("memdb")
	assert.Nil(t, err)

	err = d.Open("memdb")
	assert.Equal(t, ErrAlreadyOpen, err, "error expected when db is already open")
}
