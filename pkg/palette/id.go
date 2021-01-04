package palette

// ID is a palette id
type ID string

const (
	// EmptyID is not a palette
	EmptyID ID = ""
	// DefaultID is the default palette
	DefaultID ID = "default"
)

// IsEmptyID returns true, when the ID is empty
func IsEmptyID(id ID) bool {
	return id == EmptyID
}

// IsDefaultID returns true, when the ID is the default
func IsDefaultID(id ID) bool {
	return id == DefaultID
}
