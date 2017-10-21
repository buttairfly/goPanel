package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntMap(t *testing.T) {
	cases := []struct {
		desc                     string
		in, inMin, inMax         int
		expected, outMin, outMax int
	}{
		{
			desc: "zero",
		},
		{
			desc:     "map_to_negative_higher_scale",
			in:       4,
			inMin:    0,
			inMax:    100,
			outMin:   -100,
			outMax:   100,
			expected: -92,
		},
		{
			desc:     "map_to_positive_higher_scale",
			in:       -4,
			inMin:    -100,
			inMax:    0,
			outMin:   0,
			outMax:   100,
			expected: 96,
		},
		{
			desc:     "result gets truncated",
			in:       1,
			inMin:    0,
			inMax:    100,
			outMin:   0,
			outMax:   50,
			expected: 0,
		},
	}
	for _, c := range cases {
		out := IntMap(c.in, c.inMin, c.inMax, c.outMin, c.outMax)
		assert.Equal(t, c.expected, out, "equal mapping")
	}
}

func TestIntConstrain(t *testing.T) {
	cases := []struct {
		desc         string
		in, min, max int
		expected     int
	}{
		{
			desc: "zero",
		},
		{
			desc:     "in_in_range",
			in:       5,
			min:      0,
			max:      10,
			expected: 5,
		},
		{
			desc:     "in_at_lower",
			in:       5,
			min:      5,
			max:      10,
			expected: 5,
		},
		{
			desc:     "in_at_upper",
			in:       10,
			min:      5,
			max:      10,
			expected: 10,
		},
		{
			desc:     "in_under_lower",
			in:       4,
			min:      5,
			max:      10,
			expected: 5,
		},
		{
			desc:     "in_over_upper",
			in:       11,
			min:      5,
			max:      10,
			expected: 10,
		},
	}
	for _, c := range cases {
		out := IntConstrain(c.in, c.min, c.max)
		assert.Equal(t, c.expected, out, "equal constrain")
	}
}
