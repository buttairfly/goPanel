package pipepart

import (
	"encoding/json"
)

// PipeParam is a struct to model and marshal a pipe param
type PipeParam struct {
	Name     string        `json:"name" yaml:"name"`
	Type     PipeParamType `json:"type" yaml:"type"`
	Value    string        `json:"value" yaml:"value"`
	Readonly bool          `json:"readonly,omitempty" yaml:"readonly,omitempty"`
}

type aliasPipeParam struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Value    string `json:"value" yaml:"value"`
	Readonly bool   `json:"readonly,omitempty" yaml:"readonly,omitempty"`
}

// MarshalJSON marshals a PipeParam to json or error
func (me *PipeParam) MarshalJSON() ([]byte, error) {
	return json.Marshal(&aliasPipeParam{
		Name:     me.Name,
		Type:     me.Type.String(),
		Value:    me.Value,
		Readonly: me.Readonly,
	})
}
