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

var testNamespace = da.Namespace([]byte("test"))

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
	t.Run("Given height is from the future", func(t *testing.T) {
		HeightFromFutureTest(t, d)
	})
}

// BasicDATest tests round trip of messages to DA and back.
func BasicDATest(t *testing.T, d da.DA) {
	msg1 := []byte("message 1")
	msg2 := []byte("message 2")

	ctx := context.TODO()
	id1, err := d.Submit(ctx, []da.Blob{msg1}, 0, testNamespace)
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)

	id2, err := d.Submit(ctx, []da.Blob{msg2}, 0, testNamespace)
	assert.NoError(t, err)
	assert.NotEmpty(t, id2)

	id3, err := d.SubmitWithOptions(ctx, []da.Blob{msg1}, 0, testNamespace, []byte("random options"))
	assert.NoError(t, err)
	assert.NotEmpty(t, id3)

	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id1, id3)

	ret, err := d.Get(ctx, id1, testNamespace)
	assert.NoError(t, err)
	assert.Equal(t, []da.Blob{msg1}, ret)

	commitment1, err := d.Commit(ctx, []da.Blob{msg1}, []byte{})
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment1)

	commitment2, err := d.Commit(ctx, []da.Blob{msg2}, []byte{})
	assert.NoError(t, err)
	assert.NotEmpty(t, commitment2)

	ids := [][]byte{id1[0], id2[0], id3[0]}
	proofs, err := d.GetProofs(ctx, ids, testNamespace)
	assert.NoError(t, err)
	assert.NotEmpty(t, proofs)
	oks, err := d.Validate(ctx, ids, proofs, testNamespace)
	assert.NoError(t, err)
	assert.NotEmpty(t, oks)
	for _, ok := range oks {
		assert.True(t, ok)
	}
}

// CheckErrors ensures that errors are handled properly by DA.
func CheckErrors(t *testing.T, d da.DA) {
	ctx := context.TODO()
	blob, err := d.Get(ctx, []da.ID{[]byte("invalid")}, testNamespace)
	assert.Error(t, err)
	assert.Empty(t, blob)
}

// GetIDsTest tests iteration over DA
func GetIDsTest(t *testing.T, d da.DA) {
	msgs := [][]byte{[]byte("msg1"), []byte("msg2"), []byte("msg3")}

	ctx := context.TODO()
	ids, err := d.Submit(ctx, msgs, 0, testNamespace)
	assert.NoError(t, err)
	assert.Len(t, ids, len(msgs))

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
		assert.NotNil(t, ret)
		assert.NotZero(t, ret.Timestamp)
		if len(ret.IDs) > 0 {
			blobs, err := d.Get(ctx, ret.IDs, testNamespace)
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
			if err != nil {
				assert.Equal(t, err.Error(), ErrTooHigh.Error())
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := uint64(1); i <= 100; i++ {
			_, err := d.Submit(ctx, [][]byte{[]byte("test")}, 0, []byte{})
			assert.NoError(t, err)
		}
	}()

	wg.Wait()
}

// HeightFromFutureTest tests the case when the given height is from the future
func HeightFromFutureTest(t *testing.T, d da.DA) {
	ctx := context.TODO()
	ret, err := d.GetIDs(ctx, 999999999, []byte{})
	assert.Error(t, err)
	assert.Nil(t, ret)
}
