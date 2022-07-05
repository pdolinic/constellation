package joinclient

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/edgelesssys/constellation/bootstrapper/internal/nodelock"
	"github.com/edgelesssys/constellation/bootstrapper/role"
	"github.com/edgelesssys/constellation/internal/cloud/metadata"
	"github.com/edgelesssys/constellation/internal/constants"
	"github.com/edgelesssys/constellation/internal/file"
	"github.com/edgelesssys/constellation/internal/grpc/atlscredentials"
	"github.com/edgelesssys/constellation/internal/grpc/dialer"
	"github.com/edgelesssys/constellation/internal/grpc/testdialer"
	activationproto "github.com/edgelesssys/constellation/joinservice/joinproto"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"
	testclock "k8s.io/utils/clock/testing"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestClient(t *testing.T) {
	someErr := errors.New("failed")
	self := metadata.InstanceMetadata{Role: role.Worker, Name: "node-1"}
	peers := []metadata.InstanceMetadata{
		{Role: role.Worker, Name: "node-2", PrivateIPs: []string{"192.0.2.8"}},
		{Role: role.ControlPlane, Name: "node-3", PrivateIPs: []string{"192.0.2.1"}},
		{Role: role.ControlPlane, Name: "node-4", PrivateIPs: []string{"192.0.2.2", "192.0.2.3"}},
	}

	testCases := map[string]struct {
		role          role.Role
		clusterJoiner *stubClusterJoiner
		disk          encryptedDisk
		nodeLock      *nodelock.Lock
		apiAnswers    []any
	}{
		"on worker: metadata self: errors occur": {
			role: role.Worker,
			apiAnswers: []any{
				selfAnswer{err: someErr},
				selfAnswer{err: someErr},
				selfAnswer{err: someErr},
				selfAnswer{instance: self},
				listAnswer{instances: peers},
				activateWorkerNodeAnswer{},
			},
			clusterJoiner: &stubClusterJoiner{},
			nodeLock:      nodelock.New(),
			disk:          &stubDisk{},
		},
		"on worker: metadata self: invalid answer": {
			role: role.Worker,
			apiAnswers: []any{
				selfAnswer{},
				selfAnswer{instance: metadata.InstanceMetadata{Role: role.Worker}},
				selfAnswer{instance: metadata.InstanceMetadata{Name: "node-1"}},
				selfAnswer{instance: self},
				listAnswer{instances: peers},
				activateWorkerNodeAnswer{},
			},
			clusterJoiner: &stubClusterJoiner{},
			nodeLock:      nodelock.New(),
			disk:          &stubDisk{},
		},
		"on worker: metadata list: errors occur": {
			role: role.Worker,
			apiAnswers: []any{
				selfAnswer{instance: self},
				listAnswer{err: someErr},
				listAnswer{err: someErr},
				listAnswer{err: someErr},
				listAnswer{instances: peers},
				activateWorkerNodeAnswer{},
			},
			clusterJoiner: &stubClusterJoiner{},
			nodeLock:      nodelock.New(),
			disk:          &stubDisk{},
		},
		"on worker: metadata list: no control plane nodes in answer": {
			role: role.Worker,
			apiAnswers: []any{
				selfAnswer{instance: self},
				listAnswer{},
				listAnswer{},
				listAnswer{},
				listAnswer{instances: peers},
				activateWorkerNodeAnswer{},
			},
			clusterJoiner: &stubClusterJoiner{},
			nodeLock:      nodelock.New(),
			disk:          &stubDisk{},
		},
		"on worker: aaas ActivateNode: errors": {
			role: role.Worker,
			apiAnswers: []any{
				selfAnswer{instance: self},
				listAnswer{instances: peers},
				activateWorkerNodeAnswer{err: someErr},
				listAnswer{instances: peers},
				activateWorkerNodeAnswer{err: someErr},
				listAnswer{instances: peers},
				activateWorkerNodeAnswer{},
			},
			clusterJoiner: &stubClusterJoiner{},
			nodeLock:      nodelock.New(),
			disk:          &stubDisk{},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			clock := testclock.NewFakeClock(time.Now())
			metadataAPI := newStubMetadataAPI()
			fileHandler := file.NewHandler(afero.NewMemMapFs())

			netDialer := testdialer.NewBufconnDialer()
			dialer := dialer.New(nil, nil, netDialer)

			client := &JoinClient{
				nodeLock:    tc.nodeLock,
				timeout:     30 * time.Second,
				interval:    time.Millisecond,
				dialer:      dialer,
				disk:        tc.disk,
				joiner:      tc.clusterJoiner,
				fileHandler: fileHandler,
				metadataAPI: metadataAPI,
				clock:       clock,
				log:         zaptest.NewLogger(t),
			}

			serverCreds := atlscredentials.New(nil, nil)
			activationServer := grpc.NewServer(grpc.Creds(serverCreds))
			activationAPI := newStubActivationServiceAPI()
			activationproto.RegisterAPIServer(activationServer, activationAPI)
			port := strconv.Itoa(constants.ActivationServiceNodePort)
			listener := netDialer.GetListener(net.JoinHostPort("192.0.2.3", port))
			go activationServer.Serve(listener)
			defer activationServer.GracefulStop()

			client.Start()

			for _, a := range tc.apiAnswers {
				switch a := a.(type) {
				case selfAnswer:
					metadataAPI.selfAnswerC <- a
				case listAnswer:
					metadataAPI.listAnswerC <- a
				case activateWorkerNodeAnswer:
					activationAPI.activateWorkerNodeAnswerC <- a
				}
				clock.Step(time.Second)
			}

			client.Stop()

			assert.True(tc.clusterJoiner.joinClusterCalled)
			assert.False(client.nodeLock.TryLockOnce()) // lock should be locked
		})
	}
}

func TestClientConcurrentStartStop(t *testing.T) {
	netDialer := testdialer.NewBufconnDialer()
	dialer := dialer.New(nil, nil, netDialer)
	client := &JoinClient{
		nodeLock:    nodelock.New(),
		timeout:     30 * time.Second,
		interval:    30 * time.Second,
		dialer:      dialer,
		disk:        &stubDisk{},
		joiner:      &stubClusterJoiner{},
		fileHandler: file.NewHandler(afero.NewMemMapFs()),
		metadataAPI: &stubRepeaterMetadataAPI{},
		clock:       testclock.NewFakeClock(time.Now()),
		log:         zap.NewNop(),
	}

	wg := sync.WaitGroup{}

	start := func() {
		defer wg.Done()
		client.Start()
	}

	stop := func() {
		defer wg.Done()
		client.Stop()
	}

	wg.Add(10)
	go stop()
	go start()
	go start()
	go stop()
	go stop()
	go start()
	go start()
	go stop()
	go stop()
	go start()
	wg.Wait()

	client.Stop()
}

type stubRepeaterMetadataAPI struct {
	selfInstance  metadata.InstanceMetadata
	selfErr       error
	listInstances []metadata.InstanceMetadata
	listErr       error
}

func (s *stubRepeaterMetadataAPI) Self(_ context.Context) (metadata.InstanceMetadata, error) {
	return s.selfInstance, s.selfErr
}

func (s *stubRepeaterMetadataAPI) List(_ context.Context) ([]metadata.InstanceMetadata, error) {
	return s.listInstances, s.listErr
}

type stubMetadataAPI struct {
	selfAnswerC chan selfAnswer
	listAnswerC chan listAnswer
}

func newStubMetadataAPI() *stubMetadataAPI {
	return &stubMetadataAPI{
		selfAnswerC: make(chan selfAnswer),
		listAnswerC: make(chan listAnswer),
	}
}

func (s *stubMetadataAPI) Self(_ context.Context) (metadata.InstanceMetadata, error) {
	answer := <-s.selfAnswerC
	return answer.instance, answer.err
}

func (s *stubMetadataAPI) List(_ context.Context) ([]metadata.InstanceMetadata, error) {
	answer := <-s.listAnswerC
	return answer.instances, answer.err
}

type selfAnswer struct {
	instance metadata.InstanceMetadata
	err      error
}

type listAnswer struct {
	instances []metadata.InstanceMetadata
	err       error
}

type stubActivationServiceAPI struct {
	activateWorkerNodeAnswerC       chan activateWorkerNodeAnswer
	activateControlPlaneNodeAnswerC chan activateControlPlaneNodeAnswer

	activationproto.UnimplementedAPIServer
}

func newStubActivationServiceAPI() *stubActivationServiceAPI {
	return &stubActivationServiceAPI{
		activateWorkerNodeAnswerC: make(chan activateWorkerNodeAnswer),
	}
}

func (s *stubActivationServiceAPI) ActivateWorkerNode(_ context.Context, _ *activationproto.ActivateWorkerNodeRequest,
) (*activationproto.ActivateWorkerNodeResponse, error) {
	answer := <-s.activateWorkerNodeAnswerC
	if answer.resp == nil {
		answer.resp = &activationproto.ActivateWorkerNodeResponse{}
	}
	return answer.resp, answer.err
}

func (s *stubActivationServiceAPI) ActivateControlPlaneNode(_ context.Context, _ *activationproto.ActivateControlPlaneNodeRequest,
) (*activationproto.ActivateControlPlaneNodeResponse, error) {
	answer := <-s.activateControlPlaneNodeAnswerC
	if answer.resp == nil {
		answer.resp = &activationproto.ActivateControlPlaneNodeResponse{}
	}
	return answer.resp, answer.err
}

type activateWorkerNodeAnswer struct {
	resp *activationproto.ActivateWorkerNodeResponse
	err  error
}

type activateControlPlaneNodeAnswer struct {
	resp *activationproto.ActivateControlPlaneNodeResponse
	err  error
}

type stubClusterJoiner struct {
	joinClusterCalled bool
	joinClusterErr    error
}

func (j *stubClusterJoiner) JoinCluster(context.Context, *kubeadm.BootstrapTokenDiscovery, string, role.Role) error {
	j.joinClusterCalled = true
	return j.joinClusterErr
}

type stubDisk struct {
	openErr                error
	closeErr               error
	uuid                   string
	uuidErr                error
	updatePassphraseErr    error
	updatePassphraseCalled bool
}

func (d *stubDisk) Open() error {
	return d.openErr
}

func (d *stubDisk) Close() error {
	return d.closeErr
}

func (d *stubDisk) UUID() (string, error) {
	return d.uuid, d.uuidErr
}

func (d *stubDisk) UpdatePassphrase(string) error {
	d.updatePassphraseCalled = true
	return d.updatePassphraseErr
}