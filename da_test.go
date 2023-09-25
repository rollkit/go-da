package da_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rollkit/go-da"
)

func TestDummyDA(t *testing.T) {
	dummy := NewDummyDA()
	t.Run("ExecuteDA", func(t *testing.T) {
		ExecuteDATest(t, dummy)
	})
	t.Run("CheckErrors", func(t *testing.T) {
		CheckErrors(t, dummy)
	})
}

// TODO(tzdybal): how to get rid of this?!
type Blob = da.Blob
type ID = da.ID

func ExecuteDATest(t *testing.T, da da.DA) {
	msg1 := []byte("message 1")
	msg2 := []byte("message 2")

	id1, proof1, err := da.Submit([]Blob{msg1})
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, proof1)

	id2, proof2, err := da.Submit([]Blob{msg2})
	assert.NoError(t, err)
	assert.NotEmpty(t, id2)
	assert.NotEmpty(t, proof2)

	id3, proof3, err := da.Submit([]Blob{msg1})
	assert.NoError(t, err)
	assert.NotEmpty(t, id3)
	assert.NotEmpty(t, proof3)

	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id1, id3)

	ret, err := da.Get(id1)
	assert.NoError(t, err)
	assert.Equal(t, []Blob{msg1}, ret)

	commitment1, err := da.Commit([]Blob{msg1})
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment1)

	commitment2, err := da.Commit([]Blob{msg2})
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment2)

	oks, err := da.Validate(commitment1, proof1)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.True(t, ok)
	}

	oks, err = da.Validate(commitment1, proof2)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.False(t, ok)
	}

	oks, err = da.Validate(commitment2, proof1)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.False(t, ok)
	}
}

func CheckErrors(t *testing.T, da da.DA) {
	blob, err := da.Get([]ID{[]byte("invalid")})
	assert.Error(t, err)
	assert.Empty(t, blob)
}
