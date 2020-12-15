package pixelpipe

import (
	"context"
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
			expectedPipeline: NewEmptyFramePipeline(context.TODO(), pipepart.ID("empty pipeline"), logger),
		},
		{
			desc: "single pipe added",
			pipes: []PipeAdder{
				{
					pipe:      generatorpipe.DrawPipe(pipepart.ID("drawPipe"), logger, make(chan generatorpipe.DrawCommand, 0)),
					addBefore: "pipeline",
				},
			},
			expectedPipeline: NewEmptyFramePipeline(context.TODO(), pipepart.ID("pipeline"), logger),
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fp := NewEmptyFramePipeline(context.TODO(), pipepart.ID(c.desc), logger)
			assert.ObjectsAreEqual(c.expectedPipeline, fp)
		})
	}
}
