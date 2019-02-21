package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/buttairfly/goPanel/src/testhelper"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMainConfig(t *testing.T) {
	const testFolder = "testdata/"
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
					testFolder + "snake_vertical___c0_19-0_10-9.config",
					testFolder + "snake_vertical___c1_9-0_0-9.config",
				},
			},
			panelFile:    "panel.config",
			expectedFile: "composedMain.config",
			actualFile:   "actual.config",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			expectedFile := fmt.Sprintf("%s%s", testFolder, c.expectedFile)
			actualFile := fmt.Sprintf("%s%s", testFolder, c.actualFile)
			panelFile := fmt.Sprintf("%s%s", testFolder, c.panelFile)

			skip := false
			for _, tileConfigPath := range c.panelConfig.TileConfigPaths {
				if _, err := os.Stat(tileConfigPath); err != nil {
					t.Log(err.Error())
					skip = true
				}
			}
			if skip {
				testhelper.FailAndSkip(t, "Run: env TEST_RECORD=true go test ./... -run TestNewTileConfigSnakeMapFile")
			}

			if testhelper.RecordCall() {
				t.Logf("Write Panel Config to file %v", panelFile)
				require.NoError(t, c.panelConfig.writeToFile(panelFile))
			}

			genConfig, err := NewConfigFromPanelConfigPath(panelFile)
			require.NoError(t, err)
			require.NotNil(t, genConfig)

			if testhelper.RecordCall() {
				t.Logf("Write Main Composed Config to file %v", expectedFile)
				require.NoError(t, genConfig.WriteToFile(expectedFile))
			}

			readConfig, err2 := newConfigFromPath(expectedFile)

			require.NoError(t, err2)
			require.NotNil(t, readConfig)

			assert.True(t, cmp.Equal(readConfig, genConfig), "error read and generated main config are not equal")
			t.Log(cmp.Diff(readConfig, genConfig))

			assert.Equal(t, c.err, genConfig.WriteToFile(actualFile), "error occurred in actual file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)
		})
	}
}
