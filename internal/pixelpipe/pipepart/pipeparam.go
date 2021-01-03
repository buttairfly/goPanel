package pipepart

// PipeParam is a struct to model and marshal a pipe param
type PipeParam struct {
	Name  string        `json:"name" yaml:"name"`
	Type  PipeParamType `json:"type" yaml:"type"`
	Value string        `json:"value" yaml:"value"`
}
