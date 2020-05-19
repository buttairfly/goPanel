package arduinocom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcHexParityChar(t *testing.T) {
	cases := []struct {
		desc        string
		seed        byte
		bareCommand string
		expected    string
	}{
		{
			desc:     "empty_command_and_seed",
			expected: "0",
		},
		{
			desc:     "empty_command_a5_seed",
			seed:     0xa5,
			expected: "f",
		},
		{
			desc:        "1char_version_command",
			bareCommand: "V",
			seed:        0xa5,
			expected:    "c",
		},
		{
			desc:        "1char_latch_command",
			bareCommand: "L",
			seed:        0xa5,
			expected:    "7",
		},
		{
			desc:        "5char_init_command_0",
			bareCommand: "I0000",
			seed:        0xa5,
			expected:    "2",
		},
		{
			desc:        "5char_init_command_200",
			bareCommand: "I00c8",
			seed:        0xa5,
			expected:    "c",
		},
		{
			desc:        "5char_init_command_200_without_seed",
			bareCommand: "I00c8",
			seed:        0,
			expected:    "3",
		},
		{
			desc:        "bchar_set_color_command_200_yellow",
			bareCommand: "S00c8ff6a00",
			seed:        0xa5,
			expected:    "5",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert.Equal(t, c.expected, calcHexParityChar(c.bareCommand, c.seed))
		})
	}
}

func TestCalcHexParity(t *testing.T) {
	cases := []struct {
		desc        string
		seed        byte
		bareCommand string
		expected    string
	}{
		{
			desc:     "empty_command_and_seed",
			expected: "0\n",
		},
		{
			desc:     "empty_command_a5_seed",
			seed:     0xa5,
			expected: "f\n",
		},
		{
			desc:        "1char_version_command",
			bareCommand: "V",
			seed:        0xa5,
			expected:    "Vc\n",
		},
		{
			desc:        "5char_init_command_0",
			bareCommand: "I0000",
			seed:        0xa5,
			expected:    "I00002\n",
		},
		{
			desc:        "5char_init_command_200",
			bareCommand: "I00c8",
			seed:        0xa5,
			expected:    "I00c8c\n",
		},
		{
			desc:        "5char_init_command_200_without_seed",
			bareCommand: "I00c8",
			seed:        0,
			expected:    "I00c83\n",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert.Equal(t, c.expected, CalcHexParity(c.bareCommand, c.seed))
		})
	}
}
