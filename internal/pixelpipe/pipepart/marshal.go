package pipepart

// Marshal is a struct to marshal a PixelPiper
type Marshal struct {
	ID     ID          `json:"id" yaml:"id"`
	PrevID ID          `json:"prevId" yaml:"prevId"`
	Params []PipeParam `json:"params" yaml:"params"`
}

// Marshal marshals a PixelPiper and Pipe
func (me *Pipe) Marshal() Marshal {
	return Marshal{
		ID:     me.GetID(),
		PrevID: me.GetPrevID(),
		Params: me.GetParams(),
	}
}
