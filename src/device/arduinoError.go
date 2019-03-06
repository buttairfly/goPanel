package device

import (
	"fmt"
	"strings"

	"github.com/buttairfly/goPanel/src/config"
)

type arduinoError struct {
	description *config.ArduinoErrorDescription
	param       string
	currentChar string
}

func IsArduinoError(line string) bool {
	if line[0] == 'E' {
		return true
	}
	return false
}

func NewArduinoError(serialConfig *config.ArduinoErrorConfig, line string) (*arduinoError, error) {
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
	return &arduinoError{
		description: description,
		param:       param,
		currentChar: currentChar,
	}, nil
}

func (ae *arduinoError) Error() string {
	param := ""
	currentChar := ""
	devider := ":"
	if ae.param != "" {
		intParam, err := ae.param
		param = fmt.Sprintf("%s %s: %s", devider, ae.description.Param, ae.param)
		devider = ","
	}
	if ae.currentChar != "" {
		currentChar = fmt.Sprintf("%s %s: %v", devider, ae.description.Character, ae.currentChar)
	}
	return fmt.Sprintf("%s%s%s", ae.description.Name, param, currentChar)
}
