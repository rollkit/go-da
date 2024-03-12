package proxy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rollkit/go-da/proxy-jsonrpc"
	"github.com/rollkit/go-da/test"
)

func TestProxy(t *testing.T) {
	dummy := test.NewDummyDA()
	server := proxy.NewServer("localhost", "3450", dummy)
	err := server.Start(context.TODO())
	require.NoError(t, err)
	defer func() {
		if err := server.Stop(context.TODO()); err != nil {
			require.NoError(t, err)
		}
	}()

	client, err := proxy.NewClient(context.TODO(), "http://localhost:3450", "")
	require.NoError(t, err)
	test.RunDATestSuite(t, &client.DA)
}
