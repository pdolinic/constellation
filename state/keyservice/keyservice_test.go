package keyservice

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/edgelesssys/constellation/coordinator/atls"
	"github.com/edgelesssys/constellation/coordinator/core"
	"github.com/edgelesssys/constellation/coordinator/pubapi/pubproto"
	"github.com/edgelesssys/constellation/coordinator/role"
	"github.com/edgelesssys/constellation/state/keyservice/keyproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/test/bufconn"
)

func TestRequestKeyLoop(t *testing.T) {
	defaultInstance := core.Instance{
		Name:       "test-instance",
		ProviderID: "/test/provider",
		Role:       role.Coordinator,
		IPs:        []string{"192.0.2.1"},
	}

	testCases := map[string]struct {
		server          *stubAPIServer
		expectedCalls   int
		listResponse    []core.Instance
		dontStartServer bool
	}{
		"success": {
			server:       &stubAPIServer{requestStateDiskKeyResp: &pubproto.RequestStateDiskKeyResponse{}},
			listResponse: []core.Instance{defaultInstance},
		},
		"no error if server throws an error": {
			server: &stubAPIServer{
				requestStateDiskKeyResp: &pubproto.RequestStateDiskKeyResponse{},
				requestStateDiskKeyErr:  errors.New("error"),
			},
			listResponse: []core.Instance{defaultInstance},
		},
		"no error if the server can not be reached": {
			server:          &stubAPIServer{requestStateDiskKeyResp: &pubproto.RequestStateDiskKeyResponse{}},
			listResponse:    []core.Instance{defaultInstance},
			dontStartServer: true,
		},
		"no error if no endpoint is available": {
			server: &stubAPIServer{requestStateDiskKeyResp: &pubproto.RequestStateDiskKeyResponse{}},
		},
		"works for multiple endpoints": {
			server: &stubAPIServer{requestStateDiskKeyResp: &pubproto.RequestStateDiskKeyResponse{}},
			listResponse: []core.Instance{
				defaultInstance,
				{
					Name:       "test-instance-2",
					ProviderID: "/test/provider",
					Role:       role.Coordinator,
					IPs:        []string{"192.0.2.2"},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			keyReceived := make(chan struct{}, 1)
			listener := bufconn.Listen(1)
			defer listener.Close()

			tlsConfig, err := atls.CreateAttestationServerTLSConfig(core.NewMockIssuer())
			require.NoError(err)
			s := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))
			pubproto.RegisterAPIServer(s, tc.server)

			if !tc.dontStartServer {
				go func() { require.NoError(s.Serve(listener)) }()
			}

			keyWaiter := &KeyAPI{
				metadata:    stubMetadata{listResponse: tc.listResponse},
				keyReceived: keyReceived,
				timeout:     500 * time.Millisecond,
			}

			// notify the API a key was received after 1 second
			go func() {
				time.Sleep(1 * time.Second)
				keyReceived <- struct{}{}
			}()

			err = keyWaiter.requestKeyLoop(
				"1234",
				grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
					return listener.DialContext(ctx)
				}),
			)
			assert.NoError(err)

			s.Stop()
		})
	}
}

func TestPushStateDiskKey(t *testing.T) {
	testCases := map[string]struct {
		testAPI     *KeyAPI
		request     *keyproto.PushStateDiskKeyRequest
		errExpected bool
	}{
		"success": {
			testAPI: &KeyAPI{keyReceived: make(chan struct{}, 1)},
			request: &keyproto.PushStateDiskKeyRequest{StateDiskKey: []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")},
		},
		"key already set": {
			testAPI: &KeyAPI{
				keyReceived: make(chan struct{}, 1),
				key:         []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
			},
			request:     &keyproto.PushStateDiskKeyRequest{StateDiskKey: []byte("BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB")},
			errExpected: true,
		},
		"incorrect size of pushed key": {
			testAPI:     &KeyAPI{keyReceived: make(chan struct{}, 1)},
			request:     &keyproto.PushStateDiskKeyRequest{StateDiskKey: []byte("AAAAAAAAAAAAAAAA")},
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			_, err := tc.testAPI.PushStateDiskKey(context.Background(), tc.request)
			if tc.errExpected {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(tc.request.StateDiskKey, tc.testAPI.key)
			}
		})
	}
}

func TestResetKey(t *testing.T) {
	api := New(nil, nil, time.Second)

	api.key = []byte{0x1, 0x2, 0x3}
	api.ResetKey()
	assert.Nil(t, api.key)
}

type stubAPIServer struct {
	requestStateDiskKeyResp *pubproto.RequestStateDiskKeyResponse
	requestStateDiskKeyErr  error
	pubproto.UnimplementedAPIServer
}

func (s *stubAPIServer) GetState(ctx context.Context, in *pubproto.GetStateRequest) (*pubproto.GetStateResponse, error) {
	return nil, nil
}

func (s *stubAPIServer) ActivateAsCoordinator(in *pubproto.ActivateAsCoordinatorRequest, srv pubproto.API_ActivateAsCoordinatorServer) error {
	return nil
}

func (s *stubAPIServer) ActivateAsNode(pubproto.API_ActivateAsNodeServer) error {
	return nil
}

func (s *stubAPIServer) ActivateAdditionalNodes(in *pubproto.ActivateAdditionalNodesRequest, srv pubproto.API_ActivateAdditionalNodesServer) error {
	return nil
}

func (s *stubAPIServer) JoinCluster(ctx context.Context, in *pubproto.JoinClusterRequest) (*pubproto.JoinClusterResponse, error) {
	return nil, nil
}

func (s *stubAPIServer) TriggerNodeUpdate(ctx context.Context, in *pubproto.TriggerNodeUpdateRequest) (*pubproto.TriggerNodeUpdateResponse, error) {
	return nil, nil
}

func (s *stubAPIServer) RequestStateDiskKey(ctx context.Context, in *pubproto.RequestStateDiskKeyRequest) (*pubproto.RequestStateDiskKeyResponse, error) {
	return s.requestStateDiskKeyResp, s.requestStateDiskKeyErr
}

type stubMetadata struct {
	listResponse []core.Instance
}

func (s stubMetadata) List(ctx context.Context) ([]core.Instance, error) {
	return s.listResponse, nil
}

func (s stubMetadata) Self(ctx context.Context) (core.Instance, error) {
	return core.Instance{}, nil
}

func (s stubMetadata) GetInstance(ctx context.Context, providerID string) (core.Instance, error) {
	return core.Instance{}, nil
}

func (s stubMetadata) SignalRole(ctx context.Context, role role.Role) error {
	return nil
}

func (s stubMetadata) SetVPNIP(ctx context.Context, vpnIP string) error {
	return nil
}

func (s stubMetadata) Supported() bool {
	return true
}