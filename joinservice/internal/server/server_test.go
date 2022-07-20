package server

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"
	"time"

	attestationtypes "github.com/edgelesssys/constellation/internal/attestation/types"
	"github.com/edgelesssys/constellation/internal/constants"
	"github.com/edgelesssys/constellation/internal/file"
	"github.com/edgelesssys/constellation/internal/logger"
	"github.com/edgelesssys/constellation/joinservice/joinproto"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	kubeadmv1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestIssueJoinTicket(t *testing.T) {
	someErr := errors.New("error")
	testKey := []byte{0x1, 0x2, 0x3}
	testCert := []byte{0x4, 0x5, 0x6}
	testID := attestationtypes.ID{
		Owner:   []byte{0x4, 0x5, 0x6},
		Cluster: []byte{0x7, 0x8, 0x9},
	}
	testJoinToken := &kubeadmv1.BootstrapTokenDiscovery{
		APIServerEndpoint: "192.0.2.1",
		CACertHashes:      []string{"hash"},
		Token:             "token",
	}
	testK8sVersion := "1.23.6"

	testCases := map[string]struct {
		isControlPlane bool
		kubeadm        stubTokenGetter
		kms            stubKeyGetter
		ca             stubCA
		id             []byte
		wantErr        bool
	}{
		"worker node": {
			kubeadm: stubTokenGetter{token: testJoinToken},
			kms:     stubKeyGetter{dataKey: testKey},
			ca:      stubCA{cert: testCert},
			id:      mustMarshalID(testID),
		},
		"GetDataKey fails": {
			kubeadm: stubTokenGetter{token: testJoinToken},
			kms:     stubKeyGetter{getDataKeyErr: someErr},
			ca:      stubCA{cert: testCert},
			id:      mustMarshalID(testID),
			wantErr: true,
		},
		"loading IDs fails": {
			kubeadm: stubTokenGetter{token: testJoinToken},
			kms:     stubKeyGetter{dataKey: testKey},
			ca:      stubCA{cert: testCert},
			id:      []byte{0x1, 0x2, 0x3},
			wantErr: true,
		},
		"no ID file": {
			kubeadm: stubTokenGetter{token: testJoinToken},
			kms:     stubKeyGetter{dataKey: testKey},
			ca:      stubCA{cert: testCert},
			wantErr: true,
		},
		"GetJoinToken fails": {
			kubeadm: stubTokenGetter{getJoinTokenErr: someErr},
			kms:     stubKeyGetter{dataKey: testKey},
			ca:      stubCA{cert: testCert},
			id:      mustMarshalID(testID),
			wantErr: true,
		},
		"GetCertificate fails": {
			kubeadm: stubTokenGetter{token: testJoinToken},
			kms:     stubKeyGetter{dataKey: testKey},
			ca:      stubCA{getCertErr: someErr},
			id:      mustMarshalID(testID),
			wantErr: true,
		},
		"control plane": {
			isControlPlane: true,
			kubeadm: stubTokenGetter{
				token: testJoinToken,
				files: map[string][]byte{"test": {0x1, 0x2, 0x3}},
			},
			kms: stubKeyGetter{dataKey: testKey},
			ca:  stubCA{cert: testCert},
			id:  mustMarshalID(testID),
		},
		"GetControlPlaneCertificateKey fails": {
			isControlPlane: true,
			kubeadm:        stubTokenGetter{token: testJoinToken, certificateKeyErr: someErr},
			kms:            stubKeyGetter{dataKey: testKey},
			ca:             stubCA{cert: testCert},
			id:             mustMarshalID(testID),
			wantErr:        true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			file := file.NewHandler(afero.NewMemMapFs())
			if len(tc.id) > 0 {
				require.NoError(file.Write(filepath.Join(constants.ServiceBasePath, constants.IDFilename), tc.id, 0o644))
			}

			// IssueJoinTicket tries to read the k8s-version ConfigMap from a mounted file.
			require.NoError(file.Write(filepath.Join(constants.ServiceBasePath, constants.K8sVersion), []byte(testK8sVersion), 0o644))

			api := New(
				logger.NewTest(t),
				file,
				tc.ca,
				tc.kubeadm,
				tc.kms,
			)

			req := &joinproto.IssueJoinTicketRequest{
				DiskUuid:       "uuid",
				IsControlPlane: tc.isControlPlane,
			}
			resp, err := api.IssueJoinTicket(context.Background(), req)
			if tc.wantErr {
				assert.Error(err)
				return
			}

			var expectedIDs attestationtypes.ID
			require.NoError(json.Unmarshal(tc.id, &expectedIDs))

			require.NoError(err)
			assert.Equal(tc.kms.dataKey, resp.StateDiskKey)
			assert.Equal(expectedIDs.Cluster, resp.ClusterId)
			assert.Equal(expectedIDs.Owner, resp.OwnerId)
			assert.Equal(tc.kubeadm.token.APIServerEndpoint, resp.ApiServerEndpoint)
			assert.Equal(tc.kubeadm.token.CACertHashes[0], resp.DiscoveryTokenCaCertHash)
			assert.Equal(tc.kubeadm.token.Token, resp.Token)
			assert.Equal(tc.ca.cert, resp.KubeletCert)

			if tc.isControlPlane {
				assert.Len(resp.ControlPlaneFiles, len(tc.kubeadm.files))
			}
		})
	}
}

func mustMarshalID(id attestationtypes.ID) []byte {
	b, err := json.Marshal(id)
	if err != nil {
		panic(err)
	}
	return b
}

type stubTokenGetter struct {
	token             *kubeadmv1.BootstrapTokenDiscovery
	getJoinTokenErr   error
	files             map[string][]byte
	certificateKeyErr error
}

func (f stubTokenGetter) GetJoinToken(time.Duration) (*kubeadmv1.BootstrapTokenDiscovery, error) {
	return f.token, f.getJoinTokenErr
}

func (f stubTokenGetter) GetControlPlaneCertificatesAndKeys() (map[string][]byte, error) {
	return f.files, f.certificateKeyErr
}

type stubKeyGetter struct {
	dataKey       []byte
	getDataKeyErr error
}

func (f stubKeyGetter) GetDataKey(context.Context, string, int) ([]byte, error) {
	return f.dataKey, f.getDataKeyErr
}

type stubCA struct {
	cert       []byte
	getCertErr error
}

func (f stubCA) GetCertificate(csr []byte) ([]byte, error) {
	return f.cert, f.getCertErr
}