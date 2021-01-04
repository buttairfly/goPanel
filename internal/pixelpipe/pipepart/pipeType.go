package pipepart

// PipeType represents a id of a PixelPiperBasic
type PipeType string

// PipeType definitions as enum
const (
	// reserved
	DummyPipe PipeType = "dummy"
	Source    PipeType = "source"
	Sink      PipeType = "sink"
	Panel     PipeType = "panel"

	// structure
	FramePipeline    PipeType = "pipeline"
	PipeIntersection PipeType = "intersection"

	// generatorpipe
	DrawGenerator           PipeType = "drawGenerator"
	FullFrameFadeGenerator  PipeType = "fullFrameFadeGenerator"
	LastBlackFrameGenerator PipeType = "lastBlackFrameGenerator"
	RainbowGenerator        PipeType = "rainbowGenerator"
	SnakeGenerator          PipeType = "snakeGenerator"
	WhiteNoiseGenerator     PipeType = "whiteNoiseGenerator"

	// alphablender
	ClockBlender PipeType = "clockBlender"
)

// GetReservedPipeTypes returns all reserved PipeTypes
func GetReservedPipeTypes() []PipeType {
	return []PipeType{
		DummyPipe,
		Source,
		Sink,
		Panel,
	}
}

// GetPipeTypes returns all possible PipeTypes
func GetPipeTypes() []PipeType {
	return []PipeType{
		FramePipeline,
		PipeIntersection,
		DrawGenerator,
		FullFrameFadeGenerator,
		LastBlackFrameGenerator,
		RainbowGenerator,
		SnakeGenerator,
		WhiteNoiseGenerator,
		ClockBlender,
	}
}
