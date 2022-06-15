package vtpm

import (
	"errors"
	"io"
	"testing"

	"github.com/edgelesssys/constellation/internal/attestation/simulator"
	"github.com/google/go-tpm-tools/client"
	"github.com/google/go-tpm/tpm2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// simTPMNOPCloser is a wrapper for the generic TPM simulator with a NOP Close() method.
type simTPMNOPCloser struct {
	io.ReadWriteCloser
}

func (s simTPMNOPCloser) Close() error {
	return nil
}

func TestMarkNodeAsInitialized(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tpm, err := simulator.OpenSimulatedTPM()
	require.NoError(err)
	defer tpm.Close()
	pcrs, err := client.ReadAllPCRs(tpm)
	require.NoError(err)

	assert.NoError(MarkNodeAsInitialized(func() (io.ReadWriteCloser, error) {
		return &simTPMNOPCloser{tpm}, nil
	}, []byte{0x0, 0x1, 0x2, 0x3}, []byte{0x4, 0x5, 0x6, 0x7}))

	pcrsInitialized, err := client.ReadAllPCRs(tpm)
	require.NoError(err)

	for i := range pcrs {
		assert.NotEqual(pcrs[i].Pcrs[uint32(PCRIndexOwnerID)], pcrsInitialized[i].Pcrs[uint32(PCRIndexOwnerID)])
		assert.NotEqual(pcrs[i].Pcrs[uint32(PCRIndexClusterID)], pcrsInitialized[i].Pcrs[uint32(PCRIndexClusterID)])
	}
}

func TestFailOpener(t *testing.T) {
	assert := assert.New(t)

	assert.Error(MarkNodeAsInitialized(func() (io.ReadWriteCloser, error) { return nil, errors.New("failed") }, []byte{0x0, 0x1, 0x2, 0x3}, []byte{0x0, 0x1, 0x2, 0x3}))
}

func TestIsNodeInitialized(t *testing.T) {
	testCases := map[string]struct {
		pcrValueOwnerID   []byte
		pcrValueClusterID []byte
		wantInitialized   bool
		wantErr           bool
	}{
		"uninitialized PCRs results in uninitialized node": {},
		"initializing PCRs result in initialized node": {
			pcrValueOwnerID:   []byte{0x0, 0x1, 0x2, 0x3},
			pcrValueClusterID: []byte{0x4, 0x5, 0x6, 0x7},
			wantInitialized:   true,
		},
		"initializing ownerID alone fails": {
			pcrValueOwnerID: []byte{0x0, 0x1, 0x2, 0x3},
			wantErr:         true,
		},
		"initializing clusterID alone fails": {
			pcrValueClusterID: []byte{0x4, 0x5, 0x6, 0x7},
			wantErr:           true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := require.New(t)
			require := require.New(t)
			tpm, err := simulator.OpenSimulatedTPM()
			require.NoError(err)
			defer tpm.Close()
			if tc.pcrValueOwnerID != nil {
				require.NoError(tpm2.PCREvent(tpm, PCRIndexOwnerID, tc.pcrValueOwnerID))
			}
			if tc.pcrValueClusterID != nil {
				require.NoError(tpm2.PCREvent(tpm, PCRIndexClusterID, tc.pcrValueClusterID))
			}
			initialized, err := IsNodeInitialized(func() (io.ReadWriteCloser, error) {
				return &simTPMNOPCloser{tpm}, nil
			})
			if tc.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err)
			require.Equal(tc.wantInitialized, initialized)
		})
	}
}