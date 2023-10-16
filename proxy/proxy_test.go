package proxy_test

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rollkit/go-da/proxy"
	"github.com/rollkit/go-da/test"
)

func TestProxy(t *testing.T) {
	dummy := test.NewDummyDA()
	server := proxy.NewServer(dummy, grpc.Creds(insecure.NewCredentials()))
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	go func() {
		_ = server.Serve(lis)
	}()

	client := proxy.NewClient()
	err = client.Start(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)
	test.RunDATestSuite(t, client)
}
