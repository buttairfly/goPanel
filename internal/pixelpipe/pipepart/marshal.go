package pipepart

// Marshal is a struct to marshal a PixelPiper
type Marshal struct {
	Type        PipeType     `json:"type" yaml:"type"`
	ID          ID           `json:"id" yaml:"id"`
	PrevID      ID           `json:"prevId,omitempty" yaml:"prevId,omitempty"`
	Params      []PipeParam  `json:"params,omitempty" yaml:"params,omitempty"`
	LastPipeID  ID           `json:"lastPipeId,omitempty" yaml:"lastPipeId,omitempty"`
	FirstPipeID ID           `json:"firstPipeId,omitempty" yaml:"firstPipeId,omitempty"`
	PixelPipes  PipesMarshal `json:"pipes,omitempty" yaml:"pipes,omitempty"`
}

// PipesMarshal is the map of marshal
type PipesMarshal map[ID]*Marshal

// MarshalFromPixelPiperInterface returns a single Marshal
func MarshalFromPixelPiperInterface(i PixelPiper) *Marshal {
	return MarshalFromPixelPiperSinkInterface(i)
}

// MarshalFromPixelPiperSinkInterface returns a single Marshal
func MarshalFromPixelPiperSinkInterface(i PixelPiperSink) *Marshal {
	m := MarshalFromPixelPiperBaseInterface(i)
	m.PrevID = i.GetPrevID()
	return m
}

// MarshalFromPixelPiperBaseInterface returns a single basic Marshal
func MarshalFromPixelPiperBaseInterface(i PixelPiperBase) *Marshal {
	return &Marshal{
		Type:   i.GetType(),
		ID:     i.GetID(),
		Params: i.GetParams(),
	}
}
