package pipepart

// ID is an identifier for a pixelPiper elemement and used to name pixelPipes
type ID string

// EmptyID is a placeholder for an empty id
const EmptyID ID = ""

// SourceID is a placeholder for an source id
const SourceID ID = "|-"

// SinkID is a placeholder for an sink id
const SinkID ID = "-|"

// PanelID is a placeholder for the overall panel id
const PanelID ID = ">>"

// IsEmptyID returns true, when the ID is empty
func IsEmptyID(id ID) bool {
	return id == EmptyID
}

// IsPlaceholderID returns true, when the ID should not be used to create a pipe
func IsPlaceholderID(id ID) bool {
	return IsEmptyID(id) || id == SourceID || id == SinkID || id == PanelID
}
