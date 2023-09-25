package da_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"sync/atomic"

	"github.com/rollkit/go-da"
)

// DummyDA is a simple implementation of in-memory DA. Not production ready! Intended only for testing!
//
// Data is stored in a map, where key is a serialized sequence number. This key is returned as ID.
// Commitments are simply hashes, and proofs are ED25519 signatures.
type DummyDA struct {
	data    map[string][]byte
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
	cnt     uint64
}

func NewDummyDA() *DummyDA {
	da := &DummyDA{
		data: make(map[string][]byte),
	}
	da.pubKey, da.privKey, _ = ed25519.GenerateKey(rand.Reader)
	return da
}

var _ da.DA = &DummyDA{}

func (d *DummyDA) Get(ids []da.ID) ([]da.Blob, error) {
	blobs := make([]da.Blob, len(ids))
	for i, id := range ids {
		blob, ok := d.data[string(id)]
		if !ok {
			return nil, errors.New("no blob for given ID")
		}
		blobs[i] = blob
	}
	return blobs, nil
}

func (d *DummyDA) Commit(blobs []da.Blob) ([]da.Commitment, error) {
	commits := make([]da.Commitment, len(blobs))
	for i, blob := range blobs {
		commits[i] = d.getHash(blob)
	}
	return commits, nil
}

func (d *DummyDA) Submit(blobs []da.Blob) ([]da.ID, []da.Proof, error) {
	ids := make([]da.ID, len(blobs))
	proofs := make([]da.Proof, len(blobs))
	for i, blob := range blobs {
		id := d.nextID()
		ids[i] = id
		proofs[i] = d.getProof(id, blob)

		d.data[string(id)] = blob
	}

	return ids, proofs, nil
}

func (d *DummyDA) Validate(ids []da.ID, proofs []da.Proof) ([]bool, error) {
	if len(ids) != len(proofs) {
		return nil, errors.New("number of IDs doesn't equal to number of proofs")
	}
	results := make([]bool, len(ids))
	for i := 0; i < len(ids); i++ {
		results[i] = ed25519.Verify(d.pubKey, ids[i], proofs[i])
	}
	return results, nil
}

func (d *DummyDA) nextID() []byte {
	cnt := atomic.AddUint64(&d.cnt, 1)
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
