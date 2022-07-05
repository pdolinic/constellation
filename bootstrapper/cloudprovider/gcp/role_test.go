package gcp

import (
	"testing"

	"github.com/edgelesssys/constellation/bootstrapper/role"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		// https://github.com/census-instrumentation/opencensus-go/issues/1262
		goleak.IgnoreTopFunction("go.opencensus.io/stats/view.(*worker).start"),
	)
}

func TestExtractRole(t *testing.T) {
	testCases := map[string]struct {
		metadata map[string]string
		wantRole role.Role
	}{
		"bootstrapper role": {
			metadata: map[string]string{
				roleMetadataKey: role.ControlPlane.String(),
			},
			wantRole: role.ControlPlane,
		},
		"node role": {
			metadata: map[string]string{
				roleMetadataKey: role.Worker.String(),
			},
			wantRole: role.Worker,
		},
		"unknown role": {
			metadata: map[string]string{
				roleMetadataKey: "some-unknown-role",
			},
			wantRole: role.Unknown,
		},
		"no role": {
			wantRole: role.Unknown,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			role := extractRole(tc.metadata)

			assert.Equal(tc.wantRole, role)
		})
	}
}