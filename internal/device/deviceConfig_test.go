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

var _ filereadwriter.Yaml = (*LedDeviceConfig)(nil)

func TestNewDeviceConfigFile(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	const testFolder = "testdata/"
	cases := []struct {
		desc            string
		ledDeviceConfig *LedDeviceConfig
		expectedFile    string
		actualFile      string
		err             error
	}{
		{
			desc: "serial",
			ledDeviceConfig: &LedDeviceConfig{
				Type: Serial,
				SerialConfig: &arduinocom.SerialConfig{
					StreamConfig: &arduinocom.StreamConfig{
						Name:        "/dev/ttyUSB0",
						Baud:        115200,
						ReadTimeout: 300 * time.Millisecond,
						Size:        8,
					},
					ReadBufferSize:     1024,
					RawFramePartNumLed: 10,
					Verbose:            false,
					VerboseArduino:     false,
					ParitySeed:         0xa5,
					InitSleepTime:      7 * time.Millisecond,
					CmdSleepTime:       5 * time.Millisecond, // 5.5ms when verbose = true
				},
			},
			expectedFile: ".config.yaml",
			actualFile:   "actual.config.yaml",
		},
		{
			desc: "print",
			ledDeviceConfig: &LedDeviceConfig{
				Type: Print,
				PrintConfig: &PrintConfig{
					FramesPerSecond: 1,
					Quiet:           false,
				},
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
				t.Logf("Write Device Config to file %+v %+v", c.ledDeviceConfig, c.ledDeviceConfig.SerialConfig)
				require.NoError(t, c.ledDeviceConfig.WriteToYamlFile(expectedFile))
			}

			readConfig, err2 := NewDeviceConfigFromPath(expectedFile, logger)
			require.NoError(t, err2)
			t.Log(cmp.Diff(readConfig, c.ledDeviceConfig))
			assert.True(t, cmp.Equal(readConfig, c.ledDeviceConfig), "error read and generated device config are not equal")
			assert.Equal(t, c.err, c.ledDeviceConfig.WriteToYamlFile(actualFile), "error occurred in file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)
		})
	}
}
