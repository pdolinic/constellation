package nodestate

import (
	"path/filepath"
	"testing"

	"github.com/edgelesssys/constellation/bootstrapper/role"
	"github.com/edgelesssys/constellation/internal/file"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestFromFile(t *testing.T) {
	testCases := map[string]struct {
		fileContents string
		wantState    *NodeState
		wantErr      bool
	}{
		"nodestate exists": {
			fileContents: `{	"Role": "ControlPlane", "OwnerID": "T3duZXJJRA==", "ClusterID": "Q2x1c3RlcklE"	}`,
			wantState: &NodeState{
				Role:      role.ControlPlane,
				OwnerID:   []byte("OwnerID"),
				ClusterID: []byte("ClusterID"),
			},
		},
		"nodestate file does not exist": {
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			fs := afero.NewMemMapFs()
			if tc.fileContents != "" {
				require.NoError(fs.MkdirAll(filepath.Dir(nodeStatePath), 0o755))
				require.NoError(afero.WriteFile(fs, nodeStatePath, []byte(tc.fileContents), 0o644))
			}
			fileHandler := file.NewHandler(fs)
			state, err := FromFile(fileHandler)
			if tc.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err)
			assert.Equal(tc.wantState, state)
		})
	}
}

func TestToFile(t *testing.T) {
	testCases := map[string]struct {
		precreateFile bool
		state         *NodeState
		wantFile      string
		wantErr       bool
	}{
		"writing works": {
			state: &NodeState{
				Role:      role.ControlPlane,
				OwnerID:   []byte("OwnerID"),
				ClusterID: []byte("ClusterID"),
			},
			wantFile: `{
	"Role": "ControlPlane",
	"OwnerID": "T3duZXJJRA==",
	"ClusterID": "Q2x1c3RlcklE"
}`,
		},
		"file exists already": {
			precreateFile: true,
			wantErr:       true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			fs := afero.NewMemMapFs()
			if tc.precreateFile {
				require.NoError(fs.MkdirAll(filepath.Dir(nodeStatePath), 0o755))
				require.NoError(afero.WriteFile(fs, nodeStatePath, []byte("pre-existing"), 0o644))
			}
			fileHandler := file.NewHandler(fs)
			err := tc.state.ToFile(fileHandler)
			if tc.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err)

			fileContents, err := afero.ReadFile(fs, nodeStatePath)
			require.NoError(err)
			assert.Equal(tc.wantFile, string(fileContents))
		})
	}
}