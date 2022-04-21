package keyservice

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/edgelesssys/constellation/coordinator/atls"
	"github.com/edgelesssys/constellation/coordinator/config"
	"github.com/edgelesssys/constellation/coordinator/core"
	"github.com/edgelesssys/constellation/coordinator/pubapi/pubproto"
	"github.com/edgelesssys/constellation/state/keyservice/keyproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// KeyAPI is the interface called by the Coordinator or an admin during restart of a node.
type KeyAPI struct {
	mux         sync.Mutex
	metadata    core.ProviderMetadata
	issuer      core.QuoteIssuer
	key         []byte
	keyReceived chan struct{}
	timeout     time.Duration
	keyproto.UnimplementedAPIServer
}

// New initializes a KeyAPI with the given parameters.
func New(issuer core.QuoteIssuer, metadata core.ProviderMetadata, timeout time.Duration) *KeyAPI {
	return &KeyAPI{
		metadata:    metadata,
		issuer:      issuer,
		keyReceived: make(chan struct{}, 1),
		timeout:     timeout,
	}
}

// PushStateDiskKeyRequest is the rpc to push state disk decryption keys to a restarting node.
func (a *KeyAPI) PushStateDiskKey(ctx context.Context, in *keyproto.PushStateDiskKeyRequest) (*keyproto.PushStateDiskKeyResponse, error) {
	a.mux.Lock()
	defer a.mux.Unlock()
	if len(a.key) != 0 {
		return nil, status.Error(codes.FailedPrecondition, "node already received a passphrase")
	}
	if len(in.StateDiskKey) != config.RNGLengthDefault {
		return nil, status.Errorf(codes.InvalidArgument, "received invalid passphrase: expected length: %d, but got: %d", config.RNGLengthDefault, len(in.StateDiskKey))
	}

	a.key = in.StateDiskKey
	a.keyReceived <- struct{}{}
	return &keyproto.PushStateDiskKeyResponse{}, nil
}

// WaitForDecryptionKey notifies the Coordinator to send a decryption key and waits until a key is received.
func (a *KeyAPI) WaitForDecryptionKey(uuid, listenAddr string) ([]byte, error) {
	if uuid == "" {
		return nil, errors.New("received no disk UUID")
	}

	tlsConfig, err := atls.CreateAttestationServerTLSConfig(a.issuer)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
	keyproto.RegisterAPIServer(server, a)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	log.Printf("Waiting for decryption key. Listening on: %s", listener.Addr().String())
	go server.Serve(listener)
	defer server.GracefulStop()

	if err := a.requestKeyLoop(uuid); err != nil {
		return nil, err
	}

	return a.key, nil
}

// ResetKey resets a previously set key.
func (a *KeyAPI) ResetKey() {
	a.key = nil
}

// requestKeyLoop continuously requests decryption keys from all available Coordinators, until the KeyAPI receives a key.
func (a *KeyAPI) requestKeyLoop(uuid string, opts ...grpc.DialOption) error {
	// we do not perform attestation, since the restarting node does not need to care about notifying the correct Coordinator
	// if an incorrect key is pushed by a malicious actor, decrypting the disk will fail, and the node will not start
	tlsClientConfig, err := atls.CreateUnverifiedClientTLSConfig()
	if err != nil {
		return err
	}

	// set up for the select statement to immediately request a key, skipping the initial delay caused by using a ticker
	firstReq := make(chan struct{}, 1)
	firstReq <- struct{}{}

	ticker := time.NewTicker(a.timeout)
	defer ticker.Stop()
	for {
		select {
		// return if a key was received
		// a key can be send by
		// - a Coordinator, after the request rpc was received
		// - by a Constellation admin, at any time this loop is running on a node during boot
		case <-a.keyReceived:
			return nil
		case <-ticker.C:
			a.requestKey(uuid, tlsClientConfig, opts...)
		case <-firstReq:
			a.requestKey(uuid, tlsClientConfig, opts...)
		}
	}
}

func (a *KeyAPI) requestKey(uuid string, tlsClientConfig *tls.Config, opts ...grpc.DialOption) {
	// list available Coordinators
	endpoints, _ := core.CoordinatorEndpoints(context.Background(), a.metadata)

	log.Printf("Sending a key request to available Coordinators: %v", endpoints)
	// notify all available Coordinators to send a key to the node
	// any errors encountered here will be ignored, and the calls retried after a timeout
	for _, endpoint := range endpoints {
		ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
		conn, err := grpc.DialContext(ctx, endpoint, append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsClientConfig)))...)
		if err == nil {
			client := pubproto.NewAPIClient(conn)
			_, _ = client.RequestStateDiskKey(ctx, &pubproto.RequestStateDiskKeyRequest{DiskUuid: uuid})
			conn.Close()
		}

		cancel()
	}
}