package filereadwriter

import (
	"io"

	"go.uber.org/zap"
)

// Yaml is the interface to write and read yaml to/frpm a file
type Yaml interface {
	FromYamlFile(path string, logger *zap.Logger) error
	FromYamlReader(r io.Reader, logger *zap.Logger) error
	WriteToYamlFile(path string) error
}
