package filereadwriter

import "io"

// Yaml is the interface to write and read yaml to/frpm a file
type Yaml interface {
	FromYamlFile(path string) error
	FromYamlReader(r io.Reader) error
	WriteToYamlFile(path string) error
}
