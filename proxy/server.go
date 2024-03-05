package proxy

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/cristalhq/jwt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-jsonrpc/auth"
	logging "github.com/ipfs/go-log/v2"
	"github.com/rollkit/go-da"
)

var log = logging.Logger("go-da")

var (
	DefaultPerms   = []auth.Permission{"public"}
	ReadPerms      = []auth.Permission{"public", "read"}
	ReadWritePerms = []auth.Permission{"public", "read", "write"}
	AllPerms       = []auth.Permission{"public", "read", "write", "admin"}
)

var AuthKey = "Authorization"

type Server struct {
	srv          *http.Server
	rpc          *jsonrpc.RPCServer
	listener     net.Listener
	authDisabled bool

	started atomic.Bool

	auth jwt.Signer
}

// JWTPayload is a utility struct for marshaling/unmarshalling
// permissions into for token signing/verifying.
type JWTPayload struct {
	Allow []auth.Permission
}

func (j *JWTPayload) MarshalBinary() (data []byte, err error) {
	return json.Marshal(j)
}

// NewTokenWithPerms generates and signs a new JWT token with the given secret
// and given permissions.
func NewTokenWithPerms(secret jwt.Signer, perms []auth.Permission) ([]byte, error) {
	p := &JWTPayload{
		Allow: perms,
	}
	return jwt.NewTokenBuilder(secret).BuildBytes(p)
}

// ExtractSignedPermissions returns the permissions granted to the token by the passed signer.
// If the token isn't signed by the signer, it will not pass verification.
func ExtractSignedPermissions(signer jwt.Signer, token string) ([]auth.Permission, error) {
	tk, err := jwt.ParseAndVerifyString(token, signer)
	if err != nil {
		return nil, err
	}
	p := new(JWTPayload)
	err = json.Unmarshal(tk.RawClaims(), p)
	if err != nil {
		return nil, err
	}
	return p.Allow, nil
}

// NewSignedJWT returns a signed JWT token with the passed permissions and signer.
func NewSignedJWT(signer jwt.Signer, permissions []auth.Permission) (string, error) {
	token, err := jwt.NewTokenBuilder(signer).Build(&JWTPayload{
		Allow: permissions,
	})
	if err != nil {
		return "", err
	}
	return token.InsecureString(), nil
}

// verifyAuth is the RPC server's auth middleware. This middleware is only
// reached if a token is provided in the header of the request, otherwise only
// methods with `read` permissions are accessible.
func (s *Server) verifyAuth(_ context.Context, token string) ([]auth.Permission, error) {
	if s.authDisabled {
		return AllPerms, nil
	}
	return ExtractSignedPermissions(s.auth, token)
}

// RegisterService registers a service onto the RPC server. All methods on the service will then be
// exposed over the RPC.
func (s *Server) RegisterService(namespace string, service interface{}, out interface{}) {
	if s.authDisabled {
		s.rpc.Register(namespace, service)
		return
	}

	auth.PermissionedProxy(AllPerms, DefaultPerms, service, getInternalStruct(out))
	s.rpc.Register(namespace, out)
}

func getInternalStruct(api interface{}) interface{} {
	return reflect.ValueOf(api).Elem().FieldByName("Internal").Addr().Interface()
}

func NewServer(address, port string, authDisabled bool, secret jwt.Signer, DA da.DA) *Server {
	rpc := jsonrpc.NewServer()
	srv := &Server{
		rpc: rpc,
		srv: &http.Server{
			Addr: address + ":" + port,
			// the amount of time allowed to read request headers. set to the default 2 seconds
			ReadHeaderTimeout: 2 * time.Second,
		},
		auth:         secret,
		authDisabled: authDisabled,
	}
	srv.srv.Handler = &auth.Handler{
		Verify: srv.verifyAuth,
		Next:   rpc.ServeHTTP,
	}
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
