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

func ExecuteDATest(t *testing.T, da da.DA) {
	msg1 := []byte("message 1")
	msg2 := []byte("message 2")

	id1, proof1, err := da.Submit(msg1)
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, proof1)

	id2, proof2, err := da.Submit(msg2)
	assert.NoError(t, err)
	assert.NotEmpty(t, id2)
	assert.NotEmpty(t, proof2)

	id3, proof3, err := da.Submit(msg1)
	assert.NoError(t, err)
	assert.NotEmpty(t, id3)
	assert.NotEmpty(t, proof3)

	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id1, id3)

	ret, err := da.Get(id1)
	assert.NoError(t, err)
	assert.Equal(t, msg1, ret)

	commitment1, err := da.Commit(msg1)
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment1)

	commitment2, err := da.Commit(msg2)
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment2)

	ok, err := da.Validate(commitment1, proof1)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = da.Validate(commitment1, proof2)
	assert.NoError(t, err)
	assert.False(t, ok)

	ok, err = da.Validate(commitment2, proof1)
	assert.NoError(t, err)
	assert.False(t, ok)
}

func CheckErrors(t *testing.T, da da.DA) {
	blob, err := da.Get([]byte("invalid"))
	assert.Error(t, err)
	assert.Empty(t, blob)
}
