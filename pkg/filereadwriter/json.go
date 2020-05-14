package filereadwriter

import "io"

// JSON is the interface to write and read json to/frpm a file
type JSON interface {
	FromJsonFile(path string) error
	FromJsonReader(r io.Reader) error
	WriteToJsonFile(path string) error
}
