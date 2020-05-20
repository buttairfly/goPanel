package config

import (
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/buttairfly/goPanel/pkg/filereadwriter"
	"github.com/buttairfly/goPanel/pkg/testhelper"
)

var _ filereadwriter.Yaml = (*MainConfig)(nil)

func TestNewMainConfig(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	gopath := os.Getenv("GOPATH")
	baseFolder := "src/github.com/buttairfly/goPanel"
	testFolder := "testdata/"
	configHardwareFolder := path.Join(gopath, baseFolder, "/internal/hardware", testFolder)
	configDeviceFolder := path.Join(gopath, baseFolder, "/internal/device", testFolder)
	configArduinocomFolder := path.Join(gopath, baseFolder, "/pkg/arduinocom", testFolder)
	cases := []struct {
		desc         string
		panelConfig  *PanelConfig
		panelFile    string
		expectedFile string
		actualFile   string
		err          error
	}{
		{
			desc: "main config",
			panelConfig: &PanelConfig{
				TileConfigPath:         configHardwareFolder,
				DeviceConfigPath:       configDeviceFolder,
				ArduinoErrorConfigPath: configArduinocomFolder,
				TileConfigFiles: []string{
					"tile.snake_vertical_c0_20-0_10-10.config.yaml",
					"tile.snake_vertical_c1_10-0_0-10.config.yaml",
				},
				DeviceConfigFile:       "device.serial.config.yaml",
				ArduinoErrorConfigFile: "device.ledpanel.arduino.error.config.yaml",
			},
			panelFile:    "main.panel.config.yaml",
			expectedFile: "main.composed.config.yaml",
			actualFile:   "actual.config.yaml",
		},
		{
			desc: "main print config",
			panelConfig: &PanelConfig{
				TileConfigPath:         configHardwareFolder,
				DeviceConfigPath:       configDeviceFolder,
				ArduinoErrorConfigPath: configArduinocomFolder,
				TileConfigFiles: []string{
					"tile.single_c0_0-0_1-1.config.yaml",
				},
				DeviceConfigFile: "device.print.config.yaml",
			},
			panelFile:    "main.panel.print.config.yaml",
			expectedFile: "main.composed.print.config.yaml",
			actualFile:   "actual.config.yaml",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			expectedFile := path.Join(testFolder, c.expectedFile)
			actualFile := path.Join(testFolder, c.actualFile)
			panelFile := path.Join(testFolder, c.panelFile)

			for _, tileConfigFile := range c.panelConfig.TileConfigFiles {
				testhelper.FileExistsOrSkip(t, path.Join(c.panelConfig.TileConfigPath, tileConfigFile))
			}
			testhelper.FileExistsOrSkip(t, path.Join(c.panelConfig.DeviceConfigPath, c.panelConfig.DeviceConfigFile))
			testhelper.FileExistsOrSkip(t, path.Join(c.panelConfig.ArduinoErrorConfigPath, c.panelConfig.ArduinoErrorConfigFile))

			if testhelper.RecordCall() {
				t.Logf("write PanelConfig to file %v", panelFile)
				require.NoError(t, c.panelConfig.WriteToYamlFile(panelFile))
			}

			genConfig, err := NewMainConfigFromPanelConfigPath(panelFile, logger)
			require.NoError(t, err)
			require.NotNil(t, genConfig)

			if testhelper.RecordCall() {
				t.Logf("write MainConfig to file %v", expectedFile)
				require.NoError(t, genConfig.WriteToYamlFile(expectedFile))
			}

			readConfig, err2 := NewMainConfigFromPath(expectedFile, logger)
			require.NoError(t, err2)
			require.NotNil(t, readConfig)

			t.Log(cmp.Diff(readConfig, genConfig))
			assert.True(t, cmp.Equal(readConfig, genConfig), "error read and generated main config are not equal")

			assert.Equal(t, c.err, genConfig.WriteToYamlFile(actualFile), "error occurred in actual file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)
		})
	}
}

func TestChangedMainConfig(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	baseFolder := "src/github.com/buttairfly/goPanel/config"
	testFolder := "testdata/"
	src := path.Join(testFolder, "main.composed.config.yaml")
	dst := path.Join(gopath, baseFolder, "main.composed.config.yaml")

	assert.NoError(t, testhelper.CopyFile(src, dst))
}
