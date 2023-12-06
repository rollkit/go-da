package proxy_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rollkit/go-da/proxy"
	"github.com/rollkit/go-da/test"
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

func TestProxy(t *testing.T) {
	dummy := test.NewDummyDA(DefaultMaxBytes)
	server := proxy.NewServer(dummy, grpc.Creds(insecure.NewCredentials()))
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() {
		_ = server.Serve(lis)
	}()

	client := proxy.NewClient()
	err = client.Start(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	test.RunDATestSuite(t, client)
}
