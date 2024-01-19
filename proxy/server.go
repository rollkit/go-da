package proxy

import (
	"context"

	"google.golang.org/grpc"

	"github.com/rollkit/go-da"
	pbda "github.com/rollkit/go-da/types/pb/da"
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

func (p *proxySrv) MaxBlobSize(ctx context.Context, request *pbda.MaxBlobSizeRequest) (*pbda.MaxBlobSizeResponse, error) {
	maxBlobSize, err := p.target.MaxBlobSize(ctx)
	return &pbda.MaxBlobSizeResponse{MaxBlobSize: maxBlobSize}, err
}

func (p *proxySrv) Get(ctx context.Context, request *pbda.GetRequest) (*pbda.GetResponse, error) {
	ids := idsPB2DA(request.Ids)
	blobs, err := p.target.Get(ctx, ids)
	return &pbda.GetResponse{Blobs: blobsDA2PB(blobs)}, err
}

func (p *proxySrv) GetIDs(ctx context.Context, request *pbda.GetIDsRequest) (*pbda.GetIDsResponse, error) {
	ids, err := p.target.GetIDs(ctx, request.Height)
	if err != nil {
		return nil, err
	}

	return &pbda.GetIDsResponse{Ids: idsDA2PB(ids)}, nil
}

func (p *proxySrv) Commit(ctx context.Context, request *pbda.CommitRequest) (*pbda.CommitResponse, error) {
	blobs := blobsPB2DA(request.Blobs)
	commits, err := p.target.Commit(ctx, blobs)
	if err != nil {
		return nil, err
	}

	return &pbda.CommitResponse{Commitments: commitsDA2PB(commits)}, nil
}

func (p *proxySrv) Submit(ctx context.Context, request *pbda.SubmitRequest) (*pbda.SubmitResponse, error) {
	blobs := blobsPB2DA(request.Blobs)

	ids, proofs, err := p.target.Submit(ctx, blobs, &da.SubmitOptions{
		GasPrice:  request.GasPrice,
		Namespace: request.Namespace.GetValue(),
	})
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
	validity, err := p.target.Validate(ctx, ids, proofs)
	if err != nil {
		return nil, err
	}
	return &pbda.ValidateResponse{Results: validity}, nil
}
