package test

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rollkit/go-da"
)

// RunDATestSuite runs all tests against given DA
func RunDATestSuite(t *testing.T, d da.DA) {
	t.Run("Basic DA test", func(t *testing.T) {
		BasicDATest(t, d)
	})
	t.Run("Get IDs and all data", func(t *testing.T) {
		GetIDsTest(t, d)
	})
	t.Run("Check Errors", func(t *testing.T) {
		CheckErrors(t, d)
	})
	t.Run("Concurrent read/write test", func(t *testing.T) {
		ConcurrentReadWriteTest(t, d)
	})
}

// BasicDATest tests round trip of messages to DA and back.
func BasicDATest(t *testing.T, d da.DA) {
	msg1 := []byte("message 1")
	msg2 := []byte("message 2")

	ctx := context.TODO()
	id1, proof1, err := d.Submit(ctx, []da.Blob{msg1}, &da.SubmitOptions{
		GasPrice:  0,
		Namespace: []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, proof1)

	id2, proof2, err := d.Submit(ctx, []da.Blob{msg2}, &da.SubmitOptions{
		GasPrice:  0,
		Namespace: []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, id2)
	assert.NotEmpty(t, proof2)

	id3, proof3, err := d.Submit(ctx, []da.Blob{msg1}, &da.SubmitOptions{
		GasPrice:  0,
		Namespace: []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, id3)
	assert.NotEmpty(t, proof3)

	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id1, id3)

	ret, err := d.Get(ctx, id1)
	assert.NoError(t, err)
	assert.Equal(t, []da.Blob{msg1}, ret)

	commitment1, err := d.Commit(ctx, []da.Blob{msg1}, []byte{})
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment1)

	commitment2, err := d.Commit(ctx, []da.Blob{msg2}, []byte{})
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment2)

	oks, err := d.Validate(ctx, id1, proof1)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.True(t, ok)
	}

	oks, err = d.Validate(ctx, id2, proof2)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.True(t, ok)
	}

	oks, err = d.Validate(ctx, id1, proof2)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.False(t, ok)
	}

	oks, err = d.Validate(ctx, id2, proof1)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.False(t, ok)
	}
}

// CheckErrors ensures that errors are handled properly by DA.
func CheckErrors(t *testing.T, d da.DA) {
	ctx := context.TODO()
	blob, err := d.Get(ctx, []da.ID{[]byte("invalid")})
	assert.Error(t, err)
	assert.Empty(t, blob)
}

// GetIDsTest tests iteration over DA
func GetIDsTest(t *testing.T, d da.DA) {
	msgs := [][]byte{[]byte("msg1"), []byte("msg2"), []byte("msg3")}

	ctx := context.TODO()
	ids, proofs, err := d.Submit(ctx, msgs, &da.SubmitOptions{
		GasPrice:  0,
		Namespace: []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
	})
	assert.NoError(t, err)
	assert.Len(t, ids, len(msgs))
	assert.Len(t, proofs, len(msgs))

	found := false
	end := time.Now().Add(1 * time.Second)

	// To Keep It Simple: we assume working with DA used exclusively for this test (mock, devnet, etc)
	// As we're the only user, we don't need to handle external data (that could be submitted in real world).
	// There is no notion of height, so we need to scan the DA to get test data back.
	for i := uint64(1); !found && !time.Now().After(end); i++ {
		ret, err := d.GetIDs(ctx, i, []byte{})
		if err != nil {
			t.Error("failed to get IDs:", err)
		}
		if len(ret) > 0 {
			blobs, err := d.Get(ctx, ret)
			assert.NoError(t, err)

			// Submit ensures atomicity of batch, so it makes sense to compare actual blobs (bodies) only when lengths
			// of slices is the same.
			if len(blobs) == len(msgs) {
				found = true
				for b := 0; b < len(blobs); b++ {
					if !bytes.Equal(blobs[b], msgs[b]) {
						found = false
					}
				}
			}
		}
	}

	assert.True(t, found)
}

// ConcurrentReadWriteTest tests the use of mutex lock in DummyDA by calling separate methods that use `d.data` and making sure there's no race conditions
func ConcurrentReadWriteTest(t *testing.T, d da.DA) {
	var wg sync.WaitGroup
	wg.Add(2)

	ctx := context.TODO()

	go func() {
		defer wg.Done()
		for i := uint64(1); i <= 100; i++ {
			_, err := d.GetIDs(ctx, i, []byte{})
			assert.NoError(t, err)
		}
	}()

	go func() {
		defer wg.Done()
		for i := uint64(1); i <= 100; i++ {
			_, _, err := d.Submit(ctx, [][]byte{[]byte("test")}, &da.SubmitOptions{
				GasPrice:  0,
				Namespace: []byte{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
			})
			assert.NoError(t, err)
		}
	}()

	wg.Wait()
}
