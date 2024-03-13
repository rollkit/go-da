package proxy

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/filecoin-project/go-jsonrpc"
	logging "github.com/ipfs/go-log/v2"

	"github.com/rollkit/go-da"
)

var log = logging.Logger("go-da")

// Server is a jsonrpc service that can serve the DA interface
type Server struct {
	srv      *http.Server
	rpc      *jsonrpc.RPCServer
	listener net.Listener

	started atomic.Bool
}

// RegisterService registers a service onto the RPC server. All methods on the service will then be
// exposed over the RPC.
func (s *Server) RegisterService(namespace string, service interface{}, out interface{}) {
	s.rpc.Register(namespace, service)
}

// NewServer accepts the host address port and the DA implementation to serve as a jsonrpc service
func NewServer(address, port string, DA da.DA) *Server {
	rpc := jsonrpc.NewServer()
	srv := &Server{
		rpc: rpc,
		srv: &http.Server{
			Addr: address + ":" + port,
		},
	}
	srv.srv.Handler = http.HandlerFunc(rpc.ServeHTTP)
	srv.RegisterService("da", DA, &API{})
	return srv
}

// Start starts the RPC Server.
func (s *Server) Start(context.Context) error {
	couldStart := s.started.CompareAndSwap(false, true)
	if !couldStart {
		log.Warn("cannot start server: already started")
		return nil
	}
	listener, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}
	s.listener = listener
	log.Infow("server started", "listening on", s.srv.Addr)
	//nolint:errcheck
	go s.srv.Serve(listener)
	return nil
}

// Stop stops the RPC Server.
func (s *Server) Stop(ctx context.Context) error {
	couldStop := s.started.CompareAndSwap(true, false)
	if !couldStop {
		log.Warn("cannot stop server: already stopped")
		return nil
	}
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	s.listener = nil
	log.Info("server stopped")
	return nil
}
