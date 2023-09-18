package da

type DA interface {
	Get(id ID) (Blob, error)
	Commit(blob Blob) (Commitment, error)
	Submit(blob Blob) (ID, Proof, error)
	Validate(commit Commitment, proof Proof) (bool, error)
}

type BatchDA interface {
	Get(ids []ID) ([]Blob, error)
	Commit(blobs []Blob) ([]Commitment, error)
	Submit(blobs []Blob) ([]ID, []Proof, error)
	Validate(commits []Commitment, proofs []Proof) ([]bool, error)
}

type Blob = []byte
type ID = []byte
type Proof = []byte
type Commitment = []byte
