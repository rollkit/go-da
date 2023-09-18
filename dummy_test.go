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

func (d *DummyDA) Get(id da.ID) (da.Blob, error) {
	blob, ok := d.data[string(id)]
	if !ok {
		return nil, errors.New("no blob for given ID")
	}
	return blob, nil
}

func (d *DummyDA) Commit(blob da.Blob) (da.Commitment, error) {
	return d.getHash(blob), nil
}

func (d *DummyDA) Submit(blob da.Blob) (da.ID, da.Proof, error) {
	id := d.nextID()
	proof := d.getProof(id, blob)

	d.data[string(id)] = blob
	return id, proof, nil
}

func (d *DummyDA) Validate(commit da.Commitment, proof da.Proof) (bool, error) {
	return ed25519.Verify(d.pubKey, commit, proof), nil
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
