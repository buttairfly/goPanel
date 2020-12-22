package pixelpipe

import (
	"testing"

	"github.com/buttairfly/goPanel/internal/pixelpipe/generatorpipe"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewEmptyFramePipeline_AddPipeAfter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	type PipeAdder struct {
		pipe      pipepart.PixelPiper
		addBefore pipepart.ID
	}

	const testFolder = "testdata/"
	cases := []struct {
		desc             string
		pipes            []PipeAdder
		expectedPipeline *FramePipeline
	}{
		{
			desc:             "empty pipeline",
			expectedPipeline: NewEmptyFramePipeline("empty pipeline", logger),
		},
		{
			desc: "single pipe added",
			pipes: []PipeAdder{
				{
					pipe:      generatorpipe.DrawGenerator("drawGenerator", logger, make(chan generatorpipe.DrawCommand, 0)),
					addBefore: "pipeline",
				},
			},
			expectedPipeline: NewEmptyFramePipeline("pipeline", logger),
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fp := NewEmptyFramePipeline(pipepart.ID(c.desc), logger)
			assert.ObjectsAreEqual(c.expectedPipeline, fp)
		})
	}
}
