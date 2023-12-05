package proxy

import (
	"context"

	"google.golang.org/grpc"

	"github.com/rollkit/go-da"
	pbda "github.com/rollkit/go-da/types/pb/da"
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

// NewServer creates new gRPC Server configured to serve DA proxy.
func NewServer(d da.DA, opts ...grpc.ServerOption) *grpc.Server {
	srv := grpc.NewServer(opts...)

	proxy := &proxySrv{target: d}

	pbda.RegisterDAServiceServer(srv, proxy)

	return srv
}

type proxySrv struct {
	target da.DA
}

func (p *proxySrv) Config(ctx context.Context, request *pbda.ConfigRequest) (*pbda.ConfigResponse, error) {
	return &pbda.ConfigResponse{MaxBlobSize: DefaultMaxBytes}, nil
}

func (p *proxySrv) Get(ctx context.Context, request *pbda.GetRequest) (*pbda.GetResponse, error) {
	ids := idsPB2DA(request.Ids)
	blobs, err := p.target.Get(ids)
	return &pbda.GetResponse{Blobs: blobsDA2PB(blobs)}, err
}

func (p *proxySrv) GetIDs(ctx context.Context, request *pbda.GetIDsRequest) (*pbda.GetIDsResponse, error) {
	ids, err := p.target.GetIDs(request.Height)
	if err != nil {
		return nil, err
	}

	return &pbda.GetIDsResponse{Ids: idsDA2PB(ids)}, nil
}

func (p *proxySrv) Commit(ctx context.Context, request *pbda.CommitRequest) (*pbda.CommitResponse, error) {
	blobs := blobsPB2DA(request.Blobs)
	commits, err := p.target.Commit(blobs)
	if err != nil {
		return nil, err
	}

	return &pbda.CommitResponse{Commitments: commitsDA2PB(commits)}, nil
}

func (p *proxySrv) Submit(ctx context.Context, request *pbda.SubmitRequest) (*pbda.SubmitResponse, error) {
	blobs := blobsPB2DA(request.Blobs)

	ids, proofs, err := p.target.Submit(blobs)
	if err != nil {
		return nil, err
	}

	resp := &pbda.SubmitResponse{
		Ids:    make([]*pbda.ID, len(ids)),
		Proofs: make([]*pbda.Proof, len(proofs)),
	}

	for i := range ids {
		resp.Ids[i] = &pbda.ID{Value: ids[i]}
		resp.Proofs[i] = &pbda.Proof{Value: proofs[i]}
	}

	return resp, nil
}

func (p *proxySrv) Validate(ctx context.Context, request *pbda.ValidateRequest) (*pbda.ValidateResponse, error) {
	ids := idsPB2DA(request.Ids)
	proofs := proofsPB2DA(request.Proofs)
	//TODO implement me
	validity, err := p.target.Validate(ids, proofs)
	if err != nil {
		return nil, err
	}
	return &pbda.ValidateResponse{Results: validity}, nil
}
