package pipepart

// PipeParamType is a type of param for a pipe
type PipeParamType int

const (
	// ColorRGB is a rgb color
	ColorRGB PipeParamType = iota
	// ColorRGBA is a rgb color with alpha value
	ColorRGBA
	// ColorA is a color only with alpha value
	ColorA

	// Gauge0to1 is a float64 value from 0.0 to 1.0
	Gauge0to1

	// Float64 is a float64 value from -Inf to +Inf
	Float64

	// Integer is a int value from -int to +int
	Integer

	// UInteger is a int value from 0 to +int
	UInteger

	// NameID is a string id
	NameID
)

func (me PipeParamType) String() string {
	return [...]string{"colorRGB", "colorRGBA", "colorA", "gauge0to1", "float64", "integer", "uInteger", "nameId"}[me]
}
