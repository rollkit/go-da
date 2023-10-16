package proxy

import (
	"context"
	"google.golang.org/grpc"

	"github.com/rollkit/go-da"
	pbda "github.com/rollkit/go-da/types/pb/da"
)

type Client struct {
	conn *grpc.ClientConn

	client pbda.DAServiceClient
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Start(target string, opts ...grpc.DialOption) (err error) {
	c.conn, err = grpc.Dial(target, opts...)
	if err != nil {
		return err
	}
	c.client = pbda.NewDAServiceClient(c.conn)

	return nil
}

func (c *Client) Stop() error {
	return c.conn.Close()
}

func (c *Client) Get(ids []da.ID) ([]da.Blob, error) {
	req := &pbda.GetRequest{
		Ids: make([]*pbda.ID, len(ids)),
	}
	for i := range ids {
		req.Ids[i] = &pbda.ID{Value: ids[i]}
	}
	resp, err := c.client.Get(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return blobsPB2DA(resp.Blobs), nil
}

func (c *Client) GetIDs(height uint64) ([]da.ID, error) {
	req := &pbda.GetIDsRequest{Height: height}
	resp, err := c.client.GetIDs(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return idsPB2DA(resp.Ids), nil
}

func (c *Client) Commit(blobs []da.Blob) ([]da.Commitment, error) {
	req := &pbda.CommitRequest{
		Blobs: blobsDA2PB(blobs),
	}

	resp, err := c.client.Commit(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	return commitsPB2DA(resp.Commitments), nil
}

func (c *Client) Submit(blobs []da.Blob) ([]da.ID, []da.Proof, error) {
	req := &pbda.SubmitRequest{
		Blobs: blobsDA2PB(blobs),
	}

	resp, err := c.client.Submit(context.TODO(), req)
	if err != nil {
		return nil, nil, err
	}

	ids := make([]da.ID, len(resp.Ids))
	proofs := make([]da.Proof, len(resp.Proofs))
	for i := range resp.Ids {
		ids[i] = resp.Ids[i].Value
		proofs[i] = resp.Proofs[i].Value
	}

	return ids, proofs, nil
}

func (c *Client) Validate(ids []da.ID, proofs []da.Proof) ([]bool, error) {
	req := &pbda.ValidateRequest{
		Ids:    idsDA2PB(ids),
		Proofs: proofsDA2PB(proofs),
	}
	resp, err := c.client.Validate(context.TODO(), req)
	return resp.Results, err
}
