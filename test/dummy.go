package test

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/rollkit/go-da"
)

// The following consts are copied from appconsts to avoid dependency hell
const (
	// NamespaceVersionSize is the size of a namespace version in bytes.
	NamespaceVersionSize = 1

	// NamespaceIDSize is the size of a namespace ID in bytes.
	NamespaceIDSize = 28

	// NamespaceSize is the size of a namespace (version + ID) in bytes.
	NamespaceSize = NamespaceVersionSize + NamespaceIDSize

	// ShareSize is the size of a share in bytes.
	ShareSize = 512

	// ShareInfoBytes is the number of bytes reserved for information. The info
	// byte contains the share version and a sequence start idicator.
	ShareInfoBytes = 1

	// ContinuationSparseShareContentSize is the number of bytes usable for data
	// in a continuation sparse share of a sequence.
	ContinuationSparseShareContentSize = ShareSize - NamespaceSize - ShareInfoBytes

	// DefaultGovMaxSquareSize is the default value for the governance modifiable
	// max square size.
	DefaultGovMaxSquareSize = 64

	DefaultMaxBytes = DefaultGovMaxSquareSize * DefaultGovMaxSquareSize * ContinuationSparseShareContentSize
)

// DummyDA is a simple implementation of in-memory DA. Not production ready! Intended only for testing!
//
// Data is stored in a map, where key is a serialized sequence number. This key is returned as ID.
// Commitments are simply hashes, and proofs are ED25519 signatures.
type DummyDA struct {
	mu          *sync.Mutex // protects data and height
	data        map[uint64][]kvp
	maxBlobSize uint64
	height      uint64
	privKey     ed25519.PrivateKey
	pubKey      ed25519.PublicKey
}

type kvp struct {
	key, value []byte
}

// NewDummyDA create new instance of DummyDA
func NewDummyDA(opts ...func(*DummyDA) *DummyDA) *DummyDA {
	da := &DummyDA{
		mu:          new(sync.Mutex),
		data:        make(map[uint64][]kvp),
		maxBlobSize: DefaultMaxBytes,
	}
	for _, f := range opts {
		da = f(da)
	}
	da.pubKey, da.privKey, _ = ed25519.GenerateKey(rand.Reader)
	return da
}

var _ da.DA = &DummyDA{}

// MaxBlobSize returns the max blob size in bytes.
func (d *DummyDA) MaxBlobSize() (uint64, error) {
	return d.maxBlobSize, nil
}

// Get returns Blobs for given IDs.
func (d *DummyDA) Get(ids []da.ID) ([]da.Blob, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	blobs := make([]da.Blob, len(ids))
	for i, id := range ids {
		if len(id) < 8 {
			return nil, errors.New("invalid ID")
		}
		height := binary.LittleEndian.Uint64(id)
		found := false
		for j := 0; !found && j < len(d.data[height]); j++ {
			if bytes.Equal(d.data[height][j].key, id) {
				blobs[i] = d.data[height][j].value
				found = true
			}
		}
		if !found {
			return nil, errors.New("no blob for given ID")
		}
	}
	return blobs, nil
}

// GetIDs returns IDs of Blobs at given DA height.
func (d *DummyDA) GetIDs(height uint64) ([]da.ID, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	kvps := d.data[height]
	ids := make([]da.ID, len(kvps))
	for i, kv := range kvps {
		ids[i] = kv.key
	}
	return ids, nil
}

// Commit returns cryptographic Commitments for given blobs.
func (d *DummyDA) Commit(blobs []da.Blob) ([]da.Commitment, error) {
	commits := make([]da.Commitment, len(blobs))
	for i, blob := range blobs {
		commits[i] = d.getHash(blob)
	}
	return commits, nil
}

// Submit stores blobs in DA layer.
func (d *DummyDA) Submit(blobs []da.Blob) ([]da.ID, []da.Proof, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ids := make([]da.ID, len(blobs))
	proofs := make([]da.Proof, len(blobs))
	d.height += 1
	for i, blob := range blobs {
		ids[i] = append(d.nextID(), d.getHash(blob)...)
		proofs[i] = d.getProof(ids[i], blob)

		d.data[d.height] = append(d.data[d.height], kvp{ids[i], blob})
	}

	return ids, proofs, nil
}

// Validate checks the Proofs for given IDs.
func (d *DummyDA) Validate(ids []da.ID, proofs []da.Proof) ([]bool, error) {
	if len(ids) != len(proofs) {
		return nil, errors.New("number of IDs doesn't equal to number of proofs")
	}
	results := make([]bool, len(ids))
	for i := 0; i < len(ids); i++ {
		results[i] = ed25519.Verify(d.pubKey, ids[i][8:], proofs[i])
	}
	return results, nil
}

func (d *DummyDA) nextID() []byte {
	return d.getID(d.height)
}

func (d *DummyDA) getID(cnt uint64) []byte {
	id := make([]byte, 8)
	binary.LittleEndian.PutUint64(id, cnt)
	return id
}

func (d *DummyDA) getHash(blob []byte) []byte {
	sha := sha256.Sum256(blob)
	return sha[:]
}

func (d *DummyDA) getProof(id, blob []byte) []byte {
	sign, _ := d.privKey.Sign(rand.Reader, d.getHash(blob), &ed25519.Options{})
	return sign
}
