package config

import (
	"os"
	"path"
	"testing"

	"github.com/buttairfly/goPanel/pkg/testhelper"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMainConfig(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	testFolder := path.Join(gopath, "github.com/buttairfly/goPanel")
	cases := []struct {
		desc         string
		panelConfig  *PanelConfig
		panelFile    string
		expectedFile string
		actualFile   string
		err          error
	}{
		{
			desc: "main_config",
			panelConfig: &PanelConfig{
				TileConfigPaths: []string{
					"/internal/hardware/testdata/tile.snake_vertical_c0_20-0_10-10.config.json",
					"/internal/hardware/testdata/tile.snake_vertical_c1_10-0_0-10.config.json",
					/home/keks/code/go/github.com/buttairfly/goPanel/internal/hardware/testdata/tile.snake_vertical_c0_20-0_10-10.config.json

				},/home/keks/go/github.com/buttairfly/goPanel/internal/hardware/testdata/tile.snake_vertical_c0_20-0_10-10.config.json
				DeviceConfigPath:       "/internal/device/testdata/device.serial.config.json",
				ArduinoErrorConfigPath: "/internal/device/testdata/device.serial.arduino.error.config.json",
			},
			panelFile:    "main.panel.config.json",
			expectedFile: "main.composed.config.json",
			actualFile:   "actual.config.json",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			expectedFile := path.Join(testFolder, c.expectedFile)
			actualFile := path.Join(testFolder, c.actualFile)
			panelFile := path.Join(testFolder, c.panelFile)

			skip := false
			for _, tileConfigPath := range c.panelConfig.TileConfigPaths {
				if _, err := os.Stat(testFolder + tileConfigPath); err != nil {
					t.Log(err.Error())
					skip = true
				}
			}
			if _, err := os.Stat(testFolder + c.panelConfig.DeviceConfigPath); err != nil {
				t.Log(err.Error())
				skip = true
			}
			if _, err := os.Stat(testFolder + c.panelConfig.ArduinoErrorConfigPath); err != nil {
				t.Log(err.Error())
				skip = true
			}
			if skip {
				testhelper.FailAndSkip(t, "Re-Run: env TEST_RECORD=true go test ./...")
			}

			if testhelper.RecordCall() {
				t.Logf("Write Panel Config to file %v", panelFile)
				require.NoError(t, c.panelConfig.WriteToFile(panelFile))
			}

			genConfig, err := NewConfigFromPanelConfigPath(testFolder, c.panelFile)
			require.NoError(t, err)
			require.NotNil(t, genConfig)

			if testhelper.RecordCall() {
				t.Logf("Write Main Composed Config to file %v", expectedFile)
				require.NoError(t, genConfig.WriteToFile(expectedFile))
			}

			readConfig, err2 := newConfigFromPath(expectedFile)
			require.NoError(t, err2)
			require.NotNil(t, readConfig)

			t.Log(cmp.Diff(readConfig, genConfig))
			assert.True(t, cmp.Equal(readConfig, genConfig), "error read and generated main config are not equal")

			assert.Equal(t, c.err, genConfig.WriteToFile(actualFile), "error occurred in actual file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)
		})
	}
}
