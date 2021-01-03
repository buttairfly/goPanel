package pipepart

import "go.uber.org/zap"

// CheckNoPlaceholderID checks and fatals if id is empty
func CheckNoPlaceholderID(id ID, logger *zap.Logger) {
	if IsPlaceholderID(id) {
		logger.Fatal("PipeIDPlaceholderError", zap.Error(PipeIDPlaceholderError(id)))
	}
}
