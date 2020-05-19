package device

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/pkg/arduinocom"
	"github.com/buttairfly/goPanel/pkg/filereadwriter"
	"github.com/buttairfly/goPanel/pkg/testhelper"
)

var _ filereadwriter.Yaml = (*DeviceConfig)(nil)

func TestNewDeviceConfigFile(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	const testFolder = "testdata/"
	cases := []struct {
		desc         string
		deviceConfig *DeviceConfig
		expectedFile string
		actualFile   string
		err          error
	}{
		{
			desc: "serial",
			deviceConfig: &DeviceConfig{
				Type: Serial,
				SerialConfig: &arduinocom.SerialConfig{
					StreamConfig: &arduinocom.StreamConfig{
						Name:        "/dev/ttyUSB0",
						Baud:        1152000,
						ReadTimeout: 1 * time.Second,
						Size:        8,
					},
					ReadBufferSize:   1024,
					Verbose:          true,
					VerboseArduino:   false,
					ParitySeed:       0xa5,
					InitSleepTime:    20 * time.Millisecond,
					LatchSleepTime:   10500 * time.Microsecond,
					CommandSleepTime: 100 * time.Microsecond,
				},
			},
			expectedFile: ".config.yaml",
			actualFile:   "actual.config.yaml",
		},
		{
			desc: "print",
			deviceConfig: &DeviceConfig{
				Type: Print,
			},
			expectedFile: ".config.yaml",
			actualFile:   "actual.config.yaml",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			expectedFile := path.Join(testFolder, fmt.Sprintf("device.%s%s", c.desc, c.expectedFile))
			actualFile := path.Join(testFolder, fmt.Sprintf("device.%s_%s", c.desc, c.actualFile))

			testhelper.FileExistsOrSkip(t, expectedFile)

			if testhelper.RecordCall() {
				t.Logf("Write Device Config to file %+v %+v", c.deviceConfig, c.deviceConfig.SerialConfig)
				require.NoError(t, c.deviceConfig.WriteToYamlFile(expectedFile))
			}

			readConfig, err2 := NewDeviceConfigFromPath(expectedFile, logger)
			require.NoError(t, err2)
			t.Log(cmp.Diff(readConfig, c.deviceConfig))
			assert.True(t, cmp.Equal(readConfig, c.deviceConfig), "error read and generated device config are not equal")
			assert.Equal(t, c.err, c.deviceConfig.WriteToYamlFile(actualFile), "error occurred in file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)
		})
	}
}
