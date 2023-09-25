package da

// DA defines very generic interface for interaction with Data Availability layers.
type DA interface {
	// Get returns Blob for given ID, or an error.
	//
	// Error should be returned if ID is not formatted properly, there is no Blob for given ID or any other client-level
	// error occurred (dropped connection, timeout, etc).
	Get(ids []ID) ([]Blob, error)

	// Commit creates a Commitment for the given Blob.
	Commit(blobs []Blob) ([]Commitment, error)

	// Submit submits a Blob to Data Availability layer.
	//
	// This method is synchronous. Upon successful submission to Data Availability layer, it returns ID identifying blob
	// in DA and Proof of inclusion.
	Submit(blobs []Blob) ([]ID, []Proof, error)

	// Validate validates Commitment against Proof. This should be possible without retrieving Blob.
	Validate(ids []ID, proofs []Proof) ([]bool, error)
}

// Blob is the data submitted/received from DA interface.
type Blob = []byte

// ID should contain serialized data required by the implementation to find blob in Data Availability layer.
type ID = []byte

// Commitment should contain serialized cryptographic commitment to Blob value.
type Commitment = []byte

// Proof should contain serialized proof of inclusion (publication) of Blob in Data Availability layer.
type Proof = []byte
