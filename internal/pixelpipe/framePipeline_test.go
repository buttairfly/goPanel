package pixelpipe

import (
	"fmt"
	"testing"

	"github.com/buttairfly/goPanel/internal/hardware"
	"github.com/buttairfly/goPanel/internal/pixelpipe/alphablender"
	"github.com/buttairfly/goPanel/internal/pixelpipe/generatorpipe"
	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
	"github.com/buttairfly/goPanel/pkg/palette"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewEmptyFramePipeline_AddPipeAfter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	type PipeAdder struct {
		pipe           pipepart.PixelPiper
		addBefore      pipepart.ID
		expectedPrevID pipepart.ID
	}

	var _ pipepart.PixelPiper = (*FramePipeline)(nil)

	const sourceID = pipepart.SourceID

	cases := []struct {
		desc                string
		pipes               []PipeAdder
		expectedPipeline    *FramePipeline
		expectedFirstPipeID pipepart.ID
		expectedLastPipeID  pipepart.ID
	}{
		{
			desc:                "empty pipeline",
			expectedPipeline:    NewEmptyFramePipeline("empty pipeline", logger),
			expectedFirstPipeID: pipepart.ID("empty pipeline"),
			expectedLastPipeID:  pipepart.EmptyID,
		},
		{
			desc: "single pipe added",
			pipes: []PipeAdder{
				{
					pipe:           generatorpipe.DrawGenerator("drawGenerator", palette.NewPalette("empty"), logger, make(chan generatorpipe.DrawCommand, 0)),
					addBefore:      "single pipe added",
					expectedPrevID: sourceID,
				},
			},
			expectedPipeline:    NewEmptyFramePipeline("single pipe added", logger),
			expectedFirstPipeID: pipepart.ID("drawGenerator"),
			expectedLastPipeID:  pipepart.ID("drawGenerator"),
		},
		{
			desc: "two pipes added",
			pipes: []PipeAdder{
				{
					pipe:           generatorpipe.DrawGenerator("drawGenerator", palette.NewPalette("empty"), logger, make(chan generatorpipe.DrawCommand, 0)),
					addBefore:      "two pipes added",
					expectedPrevID: sourceID,
				},
				{
					pipe:           alphablender.NewClockBlender("clock", 0.0, 1.0, logger),
					addBefore:      "two pipes added",
					expectedPrevID: "drawGenerator",
				},
			},
			expectedPipeline:    NewEmptyFramePipeline("two pipes added", logger),
			expectedFirstPipeID: pipepart.ID("drawGenerator"),
			expectedLastPipeID:  pipepart.ID("clock"),
		},
		{
			desc: "two pipes added different order",
			pipes: []PipeAdder{
				{
					pipe:           alphablender.NewClockBlender("clock", 0.0, 1.0, logger),
					addBefore:      "two pipes added different order",
					expectedPrevID: "drawGenerator",
				},
				{
					pipe:           generatorpipe.DrawGenerator("drawGenerator", palette.NewPalette("empty"), logger, make(chan generatorpipe.DrawCommand, 0)),
					addBefore:      "clock",
					expectedPrevID: sourceID,
				},
			},
			expectedPipeline:    NewEmptyFramePipeline("two pipes added different order", logger),
			expectedFirstPipeID: pipepart.ID("drawGenerator"),
			expectedLastPipeID:  pipepart.ID("clock"),
		},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			fp := NewEmptyFramePipeline(pipepart.ID(c.desc), logger)
			fp.SetInput(sourceID, make(hardware.FrameSource))
			for _, pipe := range c.pipes {
				fp.AddPipeBefore(pipe.addBefore, pipe.pipe)
			}
			assert.Equal(t, c.expectedFirstPipeID, fp.firstPipeID)
			assert.Equal(t, c.expectedLastPipeID, fp.lastPipeID)
			for _, pipe := range c.pipes {
				currentPipe, err := fp.GetPipeByID(pipe.pipe.GetID())
				require.NoError(t, err)
				assert.Equal(t, pipe.expectedPrevID, currentPipe.GetPrevID(), fmt.Sprintf("error at pipe %s", currentPipe.GetID()))
			}
			assert.ObjectsAreEqual(c.expectedPipeline, fp)
		})
	}
}
