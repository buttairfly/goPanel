package pipepart

import (
	"fmt"
	"runtime"
)

func getOrigin(level int) string {
	if _, file, line, ok := runtime.Caller(level + 1); ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknownOrigin"
}

func getError(errorName string, errorString string) error {
	return fmt.Errorf("%s - %s: %s", getOrigin(2), errorName, errorString)
}

// OutputIDMismatchError is an error which occours when a pipe output id could not be found
func OutputIDMismatchError(id, lookupID ID) error {
	return getError("OutputIDMismatchError", fmt.Sprintf("could not find output with id %s with lookupID %s", id, lookupID))
}

// PipeIDMismatchError is an error which occours when a pipe id could not be found
func PipeIDMismatchError(id, lookupID ID) error {
	return getError("PipeIDMismatchError", fmt.Sprintf("could not find pipe with id %s with lookupID %s", id, lookupID))
}

// PipeIDNotUniqueError is an error which occours when a pipe is not unique
func PipeIDNotUniqueError(id ID) error {
	return getError("PipeIDNotUniqueError", fmt.Sprintf("pipe with id %s does already exist", id))
}

// PipeIDEmptyError is an error which occours when a pipe id is empty
func PipeIDEmptyError() error {
	return getError("PipeIDEmptyError", "could create Pipe with EmptyID")
}

// PipeIDPlaceholderError is an error which occours when a pipe id is a placeholderID
func PipeIDPlaceholderError(id ID) error {
	return getError("PipeIDPlaceholderError", fmt.Sprintf("could create Pipe with placeholderID %s", id))
}
