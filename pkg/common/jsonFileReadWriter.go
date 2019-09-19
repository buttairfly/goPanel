package common

import "io"

// JSONFileReadWriter is the interface to write json to a file
type JSONFileReadWriter interface {
	FromFile(path string) error
	FromReader(r io.Reader) error
	WriteToFile(path string) error
}
