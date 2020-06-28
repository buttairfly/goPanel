package pixelpipe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewEmptyFramePipeline_AddPipeAfter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	const testFolder = "testdata/"
	cases := []struct {
		desc  string
		pipes []struct {
			pipe      Pipe
			addBefore ID
		}
		expected *FramePipeline
	}{
		{
			desc:     "empty pipeline",
			expected: NewEmptyFramePipeline(context.TODO(), ID("empty pipeline"), logger),
		},
		{
			desc:     "single pipe added",
			expected: NewEmptyFramePipeline(context.TODO(), ID("empty pipeline"), logger),
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fp := NewEmptyFramePipeline(context.TODO(), ID(c.desc), logger)
			assert.ObjectsAreEqual(c.expected, fp)
		})
	}
}
