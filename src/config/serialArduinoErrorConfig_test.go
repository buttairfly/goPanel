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

func TestNewArduinoErrorConfigFile(t *testing.T) {
	const testFolder = "testdata/"
	cases := []struct {
		desc               string
		arduinoErrorConfig *ArduinoErrorConfig
		expectedFile       string
		actualFile         string
		err                error
	}{

		{
			desc: "serial.arduino.error",
			arduinoErrorConfig: &ArduinoErrorConfig{
				"Encmd": ArduinoErrorDescription{
					Name:      "no command",
					Character: "command",
				},
				"Endcm": ArduinoErrorDescription{
					Name:      "not defined command",
					Character: "command",
				},
				"Eulet": ArduinoErrorDescription{
					Name:      "unknown letter",
					Character: "letter",
					Param:     "charType (0: undefined, 1: command, 2: hexnumber, 3: linebreak (return))",
				},
				"Eucmd": ArduinoErrorDescription{
					Name:      "unknown command",
					Character: "current char",
					Param:     "command",
				},
				"Euret": ArduinoErrorDescription{
					Name:      "unknown return",
					Character: "current char",
				},
				"Enret": ArduinoErrorDescription{
					Name:      "no return",
					Character: "current char",
				},
				"Enpov": ArduinoErrorDescription{
					Name:      "number parameter overflow",
					Character: "current char",
					Param:     "number paramerter",
				},
				"Enpeq": ArduinoErrorDescription{
					Name:      "number parameter overflow equals num leds",
					Character: "current char",
					Param:     "number paramerter",
				},
				"Euini": ArduinoErrorDescription{
					Name: "no initialisation possible",
				},
				"Enini": ArduinoErrorDescription{
					Name:      "not initialized",
					Character: "current char",
				},
				"Elati": ArduinoErrorDescription{
					Name:  "latch timeout",
					Param: "min latch wait time ms",
				},
				"Enhxn": ArduinoErrorDescription{
					Name:      "not hex number parameter",
					Character: "current char",
				},
				"Enhxc": ArduinoErrorDescription{
					Name:      "not hex color parameter",
					Character: "current char",
				},
				"Enebs": ArduinoErrorDescription{
					Name:      "not enough bytes color param",
					Character: "current char",
					Param:     "current number of bytes",
				},
			},
			expectedFile: ".config.json",
			actualFile:   "actual.config.json",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			expectedFile := fmt.Sprintf("%sdevice.%s%s", testFolder, c.desc, c.expectedFile)
			actualFile := fmt.Sprintf("%sdevice.%s_%s", testFolder, c.desc, c.actualFile)

			if testhelper.RecordCall() {
				t.Logf("Write Serial Error Config to file %s %+v", expectedFile, c.arduinoErrorConfig)
				require.NoError(t, c.arduinoErrorConfig.WriteToFile(expectedFile))
			}

			readConfig, err2 := NewArduinoErrorConfigFromPath(expectedFile)
			require.NoError(t, err2)
			t.Log(cmp.Diff(readConfig, c.arduinoErrorConfig))
			assert.True(t, cmp.Equal(readConfig, c.arduinoErrorConfig), "error read and generated serial error config are not equal")
			assert.Equal(t, c.err, c.arduinoErrorConfig.WriteToFile(actualFile), "error occurred in file write")
			defer os.Remove(actualFile)
			testhelper.Diff(t, expectedFile, actualFile)
		})
	}
}
