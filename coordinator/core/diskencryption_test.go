package core

import (
	"errors"
	"testing"

	"github.com/edgelesssys/constellation/cli/file"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetDiskUUID(t *testing.T) {
	testCases := map[string]struct {
		expectedUUID string
		openErr      error
		uuidErr      error
		errExpected  bool
	}{
		"getting uuid works": {
			expectedUUID: "uuid",
		},
		"open can fail": {
			openErr:     errors.New("open-error"),
			errExpected: true,
		},
		"getting disk uuid can fail": {
			uuidErr:     errors.New("uuid-err"),
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			zapLogger, err := zap.NewDevelopment()
			require.NoError(err)
			diskStub := encryptedDiskStub{
				openErr: tc.openErr,
				uuidErr: tc.uuidErr,
				uuid:    tc.expectedUUID,
			}
			core, err := NewCore(&stubVPN{}, nil, nil, nil, nil, nil, &diskStub, zapLogger, nil, nil, file.NewHandler(afero.NewMemMapFs()))
			require.NoError(err)
			uuid, err := core.GetDiskUUID()
			if tc.errExpected {
				assert.Error(err)
				return
			}
			require.NoError(err)

			assert.Equal(tc.expectedUUID, uuid)
		})
	}
}

func TestUpdateDiskPassphrase(t *testing.T) {
	testCases := map[string]struct {
		openErr             error
		updatePassphraseErr error
		errExpected         bool
	}{
		"updating passphrase works": {},
		"open can fail": {
			openErr:     errors.New("open-error"),
			errExpected: true,
		},
		"updating disk passphrase can fail": {
			updatePassphraseErr: errors.New("update-err"),
			errExpected:         true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			zapLogger, err := zap.NewDevelopment()
			require.NoError(err)
			diskStub := encryptedDiskStub{
				openErr:             tc.openErr,
				updatePassphraseErr: tc.updatePassphraseErr,
			}
			core, err := NewCore(&stubVPN{}, nil, nil, nil, nil, nil, &diskStub, zapLogger, nil, nil, file.NewHandler(afero.NewMemMapFs()))
			require.NoError(err)
			err = core.UpdateDiskPassphrase("passphrase")
			if tc.errExpected {
				assert.Error(err)
				return
			}
			require.NoError(err)
		})
	}
}

type encryptedDiskStub struct {
	openErr             error
	closeErr            error
	uuid                string
	uuidErr             error
	updatePassphraseErr error
}

func (s *encryptedDiskStub) UUID() (string, error) {
	return s.uuid, s.uuidErr
}

func (s *encryptedDiskStub) UpdatePassphrase(passphrase string) error {
	return s.updatePassphraseErr
}

func (s *encryptedDiskStub) Open() error {
	return s.openErr
}

func (s *encryptedDiskStub) Close() error {
	return s.closeErr
}