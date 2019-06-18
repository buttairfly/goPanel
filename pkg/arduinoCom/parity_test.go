package com

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
			desc:        "2char_command_0",
			bareCommand: "I01",
			seed:        0xa5,
			expected:    "3",
		},
		{
			desc:        "2char_command_1",
			bareCommand: "I01",
			seed:        0xa5,
			expected:    "3",
		},
		{
			desc:        "4char_command_0",
			bareCommand: "I0000",
			seed:        0xa5,
			expected:    "2",
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
			desc:        "2char_command_0",
			bareCommand: "I01",
			seed:        0xa5,
			expected:    "I013\n",
		},
		{
			desc:        "2char_command_1",
			bareCommand: "I01",
			seed:        0xa5,
			expected:    "I013\n",
		},
		{
			desc:        "4char_command_0",
			bareCommand: "I0000",
			seed:        0xa5,
			expected:    "I00002\n",
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			assert.Equal(t, c.expected, CalcHexParityCommand(c.bareCommand, c.seed))
		})
	}
}
