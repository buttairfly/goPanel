package arduinocom

import (
	"fmt"
	"strings"
)

// ArduinoError is the error struct of an arduino error which implements error interface
type ArduinoError struct {
	description *ArduinoErrorDescription
	param       string
	currentChar string
}

// IsArduinoError checks wheather the line is an error and returns a boolean
func IsArduinoError(line string) bool {
	if len(line) == 0 || line[0] == 'E' {
		return true
	}
	return false
}

// NewArduinoError looks up the error code received from the serial connection and returns the readable error
func NewArduinoError(serialConfig *ArduinoErrorConfig, line string) (*ArduinoError, error) {
	if !IsArduinoError(line) {
		return nil, fmt.Errorf("No error line beginning in line: '%s'", line)
	}
	lineParts := strings.Split(line, ",")
	description, err := serialConfig.GetDescription(lineParts[0])
	if err != nil {
		return nil, err
	}
	param, currentChar := "", ""
	for i := 1; i < len(lineParts); i++ {
		linePart := lineParts[i]
		if linePart[0] == 'p' {
			param = linePart[1:]
		}
		if linePart[0] == 's' {
			currentChar = linePart[1:]
			if currentChar == "\n" {
				currentChar = "\\n"
			}
		}
	}
	return &ArduinoError{
		description: description,
		param:       param,
		currentChar: currentChar,
	}, nil
}

func (ae *ArduinoError) Error() string {
	param := ""
	currentChar := ""
	devider := ":"
	if ae.currentChar != "" {
		currentChar = fmt.Sprintf("%s %s: %v", devider, ae.description.Character, ae.currentChar)
	}
	if ae.param != "" {
		devider = ","
		param = fmt.Sprintf("%s %s: %s", devider, ae.description.Param, ae.param)
	}
	return fmt.Sprintf("%s%s%s", ae.description.Name, param, currentChar)
}
