package filereadwriter

import (
	"io"

	"go.uber.org/zap"
)

// JSON is the interface to write and read json to/frpm a config file
type JSON interface {
	FromJsonFile(filePath string, logger *zap.Logger) error
	FromJsonReader(r io.Reader, logger *zap.Logger) error
	WriteToJsonFile(filePath string) error
}
