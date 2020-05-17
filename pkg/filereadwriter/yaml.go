package filereadwriter

import (
	"io"

	"go.uber.org/zap"
)

// Yaml is the interface to write and read yaml to/frpm a config file
type Yaml interface {
	FromYamlFile(filePath string, logger *zap.Logger) error
	FromYamlReader(r io.Reader, logger *zap.Logger) error
	WriteToYamlFile(filePath string) error
}
