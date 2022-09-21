/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package resources

import (
	"testing"

	"github.com/edgelesssys/constellation/v2/internal/kubernetes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeMaintenanceOperatorMarshalUnmarshal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	nmoDepl := NewNodeMaintenanceOperatorDeployment()
	data, err := nmoDepl.Marshal()
	require.NoError(err)

	var recreated nodeMaintenanceOperatorDeployment
	require.NoError(kubernetes.UnmarshalK8SResources(data, &recreated))
	assert.Equal(nmoDepl, &recreated)
}
