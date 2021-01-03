package pipepart

// PipeType represents a id of a PixelPiperBasic
type PipeType string

// PipeType definitions as enum
const (
	Source PipeType = "source"
	Sink   PipeType = "sink"

	DummyPipe PipeType = "dummy"

	Panel            PipeType = "panel"
	FramePipeline    PipeType = "pipeline"
	PipeIntersection PipeType = "intersection"

	// generatorpipe
	DrawGenerator           PipeType = "drawGenerator"
	FullFrameFadeGenerator  PipeType = "fullFrameFadeGenerator"
	LastBlackFrameGenerator PipeType = "lastBlackFrameGenerator"
	RainbowGenerator        PipeType = "rainbowGenerator"
	SnakeGenerator          PipeType = "snakeGenerator"
	WhiteNoiseGenerator     PipeType = "whiteNoiseGenerator"

	//alphablender
	ClockBlender PipeType = "clockBlender"
)
