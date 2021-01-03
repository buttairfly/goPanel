package pipepart

// Marshal is a struct to marshal a PixelPiper
type Marshal struct {
	//Type        PipeType       `json:"type" yaml:"type"`
	ID          ID             `json:"id" yaml:"id"`
	PrevID      ID             `json:"prevId,omitempty" yaml:"prevId,omitempty"`
	Params      []PipeParam    `json:"params,omitempty" yaml:"params,omitempty"`
	LastPipeID  ID             `json:"lastPipeId,omitempty" yaml:"lastPipeId,omitempty"`
	FirstPipeID ID             `json:"firstPipeId,omitempty" yaml:"firstPipeId,omitempty"`
	PixelPipes  map[ID]Marshal `json:"pipes,omitempty" yaml:"pipes,omitempty"`
}

// Marshal marshals a PixelPiper and Pipe
func (me *Pipe) Marshal() Marshal {
	return Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}
