package proxy_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rollkit/go-da/proxy"
	"github.com/rollkit/go-da/test"
)

func TestProxy(t *testing.T) {
	dummy := test.NewDummyDA()
	server := proxy.NewServer("localhost", "3450", true, nil, dummy)
	err := server.Start(context.TODO())
	require.NoError(t, err)
	defer server.Stop(context.TODO())

	client, err := proxy.NewClient(context.TODO(), "http://localhost:3450", "")
	require.NoError(t, err)
	test.RunDATestSuite(t, &client.DA)
}
