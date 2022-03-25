package cryptmapper

import (
	"context"
	"errors"
	"testing"

	"github.com/edgelesssys/constellation/mount/kms"
	cryptsetup "github.com/martinjungblut/go-cryptsetup"
	"github.com/stretchr/testify/assert"
)

var testDEK = []byte{
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
	0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
}

type stubCryptDevice struct {
	initErr       error
	activateErr   error
	deactivateErr error
	formatErr     error
	loadErr       error
	wipeErr       error
}

func (c *stubCryptDevice) Init(devicePath string) error {
	return c.initErr
}

func (c *stubCryptDevice) ActivateByVolumeKey(deviceName, volumeKey string, volumeKeySize, flags int) error {
	return c.activateErr
}

func (c *stubCryptDevice) Deactivate(deviceName string) error {
	return c.deactivateErr
}

func (c *stubCryptDevice) Format(deviceType cryptsetup.DeviceType, genericParams cryptsetup.GenericParams) error {
	return c.formatErr
}

func (c *stubCryptDevice) Free() bool {
	return true
}

func (c *stubCryptDevice) Load(cryptsetup.DeviceType) error {
	return c.loadErr
}

func (c *stubCryptDevice) Wipe(devicePath string, pattern int, offset, length uint64, wipeBlockSize int, flags int, progress func(size, offset uint64) int) error {
	return c.wipeErr
}

func TestCloseCryptDevice(t *testing.T) {
	testCases := map[string]struct {
		mapper      *stubCryptDevice
		errExpected bool
	}{
		"success": {
			mapper:      &stubCryptDevice{},
			errExpected: false,
		},
		"error on Init": {
			mapper: &stubCryptDevice{
				initErr: errors.New("error"),
			},
			errExpected: true,
		},
		"error on Deactivate": {
			mapper: &stubCryptDevice{
				deactivateErr: errors.New("error"),
			},
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			err := closeCryptDevice(tc.mapper, "/dev/some-device", "volume0", "test")
			if tc.errExpected {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}

	mapper := New(kms.NewStaticKMS(), &stubCryptDevice{})
	err := mapper.CloseCryptDevice("volume01-unit-test")
	assert.NoError(t, err)
}

func TestOpenCryptDevice(t *testing.T) {
	someErr := errors.New("error")

	testCases := map[string]struct {
		source      string
		volumeID    string
		dek         string
		integrity   bool
		mapper      *stubCryptDevice
		diskInfo    func(disk string) (string, error)
		errExpected bool
	}{
		"success with Load": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(testDEK),
			mapper:      &stubCryptDevice{},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: false,
		},
		"success with error on Load": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(testDEK),
			mapper:      &stubCryptDevice{loadErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: false,
		},
		"success with integrity": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(append(testDEK, testDEK[:32]...)),
			integrity:   true,
			mapper:      &stubCryptDevice{loadErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: false,
		},
		"incorrect key size": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         "short",
			mapper:      &stubCryptDevice{},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: true,
		},
		"error on Init": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(testDEK),
			mapper:      &stubCryptDevice{initErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: true,
		},
		"error on Format": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(testDEK),
			mapper:      &stubCryptDevice{loadErr: someErr, formatErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: true,
		},
		"error on Activate": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(testDEK),
			mapper:      &stubCryptDevice{activateErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: true,
		},
		"error on diskInfo": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(testDEK),
			mapper:      &stubCryptDevice{loadErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", someErr },
			errExpected: true,
		},
		"disk is already formatted": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(testDEK),
			mapper:      &stubCryptDevice{loadErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "ext4", nil },
			errExpected: true,
		},
		"error with integrity on wipe": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(append(testDEK, testDEK[:32]...)),
			integrity:   true,
			mapper:      &stubCryptDevice{loadErr: someErr, wipeErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: true,
		},
		"error with integrity on activate": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(append(testDEK, testDEK[:32]...)),
			integrity:   true,
			mapper:      &stubCryptDevice{loadErr: someErr, activateErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: true,
		},
		"error with integrity on deactivate": {
			source:      "/dev/some-device",
			volumeID:    "volume0",
			dek:         string(append(testDEK, testDEK[:32]...)),
			integrity:   true,
			mapper:      &stubCryptDevice{loadErr: someErr, deactivateErr: someErr},
			diskInfo:    func(disk string) (string, error) { return "", nil },
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			out, err := openCryptDevice(tc.mapper, tc.source, tc.volumeID, tc.dek, tc.integrity, tc.diskInfo)
			if tc.errExpected {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(cryptPrefix+tc.volumeID, out)
			}
		})
	}

	mapper := New(kms.NewStaticKMS(), &stubCryptDevice{})
	_, err := mapper.OpenCryptDevice(context.Background(), "/dev/some-device", "volume01", false)
	assert.NoError(t, err)
}