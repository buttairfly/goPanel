package pixelpipe

import (
	"fmt"

	"github.com/buttairfly/goPanel/internal/pixelpipe/pipepart"
)

// Marshal implements PixelPiper interface
func (me *FramePipeline) Marshal() pipepart.Marshal {
	pp := make(map[pipepart.ID]pipepart.Marshal, len(me.pixelPipes))
	for id, p := range me.pixelPipes {
		pp[id] = p.Marshal()
	}
	return pipepart.Marshal{
		ID:          me.GetID(),
		PrevID:      me.GetPrevID(),
		FirstPipeID: me.firstPipeID,
		LastPipeID:  me.lastPipeID,
		PixelPipes:  pp,
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
