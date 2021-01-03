package pixelpipe

import (
	"fmt"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

// Marshal is a mapping struct to marshal a FramePipeline
type Marshal struct {
	PixelPipes  map[pipepart.ID]pipepart.Marshal `json:"pipes" yaml:"pipes"`
	ID          pipepart.ID                      `json:"id" yaml:"id"`
	LastPipeID  pipepart.ID                      `json:"lastPipeId" yaml:"lastPipeId"`
	FirstPipeID pipepart.ID                      `json:"firstPipeId" yaml:"firstPipeId"`
	PrevID      pipepart.ID                      `json:"prevId" yaml:"prevId"`
}

// MarshalFramePipeline converts a marshalable palette to palette.Marshal
func (me *FramePipeline) MarshalFramePipeline() Marshal {
	pp := make(map[pipepart.ID]pipepart.Marshal, len(me.pixelPipes))
	for id, p := range me.pixelPipes {
		pp[id] = p.Marshal()
	}
	return Marshal{
		PixelPipes:  pp,
		ID:          me.GetID(),
		FirstPipeID: me.firstPipeID,
		LastPipeID:  me.lastPipeID,
		PrevID:      me.prevID,
	}
}

// Marshal implements PixelPiper interface
func (me *FramePipeline) Marshal() pipepart.Marshal {
	return pipepart.Marshal{
		ID:     me.GetID(),
		PrevID: me.prevID,
	}
}

// GetPipeByID returns a Pipe
func (me *FramePipeline) GetPipeByID(id pipepart.ID) (pipepart.PixelPiper, error) {
	if id != me.GetID() {
		pipe, ok := me.pixelPipes[id]
		if ok {
			return pipe, nil
		}
		return nil, fmt.Errorf("pipe %s not found", id)
	}
	return me, nil
}
