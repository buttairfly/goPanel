package panel

import (
	"image"

	"github.com/buttairfly/goPanel/internal/config"
	"github.com/buttairfly/goPanel/pkg/marshal"
)

// ConsumerConfig is the config used by consumers to display the current ledPix layout
type ConsumerConfig struct {
	Frame      marshal.Rectangle   `json:"frame" yaml:"frame"`
	TileFrames []marshal.Rectangle `json:"tileFrames" yaml:"tileFrames"`
}

// GetMainConfig returns the parsed main config
func (p *Panel) GetMainConfig() *config.MainConfig {
	return p.config
}

// GetConsumerConfig returns the parsed main config
func (p *Panel) GetConsumerConfig() *ConsumerConfig {
	mc := p.GetMainConfig()
	frame := image.Rectangle{}
	tileFrames := make([]marshal.Rectangle, len(mc.TileConfigs))
	for i, tile := range mc.TileConfigs {
		frame = frame.Union(tile.ToTileConfig().Bounds)
		tileFrames[i] = tile.Bounds
	}

	return &ConsumerConfig{
		Frame:      marshal.FromImageRectangle(frame),
		TileFrames: tileFrames,
	}
}