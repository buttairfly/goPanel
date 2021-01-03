package pipepart

// PipeParamType is a type of param for a pipe
type PipeParamType int

const (
	// UnknownParamType will trow an error and should not be used
	UnknownParamType PipeParamType = iota
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

	// ColorRGB is a rgb color
	ColorRGB
	// ColorRGBA is a rgb color with alpha value
	ColorRGBA
	// ColorA is a color only with alpha value
	ColorA
)

func (me PipeParamType) String() string {
	return [...]string{
		"unknown",
		"gauge0to1",
		"float64",
		"integer",
		"uInteger",
		"nameId",
		"colorRGB",
		"colorRGBA",
		"colorA",
	}[me]
}

//PipeParamTypeFromString returns a PipeParamType from a string
func PipeParamTypeFromString(s string) PipeParamType {
	switch s {
	case "gauge0to1":
		return Gauge0to1
	case "float64":
		return Float64
	case "integer":
		return Integer
	case "uInteger":
		return UInteger
	case "nameId":
		return NameID
	case "colorRGB":
		return ColorRGB
	case "colorRGBA":
		return ColorRGBA
	case "colorA":
		return ColorA
	default:
		return UnknownParamType
	}
}
