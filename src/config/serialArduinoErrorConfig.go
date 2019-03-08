package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
)

// ArduinoErrorConfig is the arduino config error map
type ArduinoErrorConfig map[string]ArduinoErrorDescription

// ArduinoErrorDescription is the description of an single arduino error and
// its parameters
type ArduinoErrorDescription struct {
	Name      string `json:"name"`
	Param     string `json:"param,omitempty"`
	Character string `json:"character,omitempty"`
}

// NewArduinoErrorConfigFromPath reads a ArduinoErrorConfig from file
func NewArduinoErrorConfigFromPath(path string) (*ArduinoErrorConfig, error) {
	aec := new(ArduinoErrorConfig)
	err := aec.FromFile(path)
	if err != nil {
		return nil, err
	}
	return aec, nil
}

// FromFile reads the config from a file at path
func (aec *ArduinoErrorConfig) FromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not read Config file %v. error: %v", path, err)
	}
	defer f.Close()
	return aec.FromReader(f)
}

// FromReader decodes the config from io.Reader
func (aec *ArduinoErrorConfig) FromReader(r io.Reader) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&*aec)
	if err != nil {
		return fmt.Errorf("can not decode json. error: %v", err)
	}
	return nil
}

// WriteToFile writes the config to a file at path
func (aec *ArduinoErrorConfig) WriteToFile(path string) error {
	jsonConfig, err := json.MarshalIndent(aec, "", "\t")
	if err != nil {
		return err
	}
	jsonConfig = append(jsonConfig, byte('\n'))
	return ioutil.WriteFile(path, jsonConfig, 0622)
}

// GetDescription returns the ArduinoErrorDescription for the key or an error on unknown key
func (aec *ArduinoErrorConfig) GetDescription(errorCode string) (*ArduinoErrorDescription, error) {
	description, ok := (*aec)[errorCode]
	if !ok {
		return nil, fmt.Errorf("Unkown error code %s", errorCode)
	}
	return &description, nil
}

// ToCppFile returns a []byte string to write to a file
func (aec *ArduinoErrorConfig) ToCppFile(filePath, name string) error {
	_, thisFile, _, _ := runtime.Caller(0)
	_, thisFileName := path.Split(thisFile)
	_, fileName := path.Split(filePath)
	_, name = path.Split(name)
	if fileName != name {
		fileName = name
	}
	upperFileName := replaceAllDots(strings.ToUpper(fileName))
	header := "// THIS FILE IS AUTO-GENERATED - DO NOT EDIT MANUALLY\n"
	header += fmt.Sprintf("// It is generated by %s\n\n", thisFileName)
	header += fmt.Sprintf("// %s\n", fileName)
	header += fmt.Sprintf("#ifndef %s\n", upperFileName)
	header += fmt.Sprintf("#define %s\n", upperFileName)
	header += "\n"
	header += "#include \"Arduino.h\"\n"
	footer := fmt.Sprintf("#endif // %s\n", fileName)

	defines := ""
	keySlice := make([]string, len(*aec), len(*aec))
	i := 0
	for errorKey := range *aec {
		keySlice[i] = errorKey
		i++
	}
	sort.Sort(sort.StringSlice(keySlice))
	for _, errorKey := range keySlice {
		errorInfo := (*aec)[errorKey]
		errorVarName := errorInfo.Name
		errorComment := ""
		commentStart := " // "
		if errorInfo.Param != "" {
			errorComment += fmt.Sprintf("%sparam: %s", commentStart, errorInfo.Param)
		}
		if errorInfo.Character != "" {
			if errorComment != "" {
				commentStart = ", "
			}
			errorComment += fmt.Sprintf("%scharacter: %s", commentStart, errorInfo.Character)
		}
		defines += fmt.Sprintf("const String Error%s = \"%s\";%s\n", ToCppVarName(errorVarName), errorKey, errorComment)
	}

	content := []byte(fmt.Sprintf("%s\n%s\n%s", header, defines, footer))
	return ioutil.WriteFile(filePath, content, 0644)
}

// ToCppVarName converts a string with spaces into a valid variable name
func ToCppVarName(s string) string {
	s = strings.Replace(s, ".", " ", -1)
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Replace(s, "-", " ", -1)
	s = strings.Replace(s, "\t", " ", -1)
	varParts := strings.Split(s, " ")
	for i, part := range varParts {
		if len(part) > 0 {
			partEnd := ""
			if len(part) > 1 {
				partEnd = part[1:]
			}
			varParts[i] = strings.ToUpper(string(part[0])) + partEnd
		}
	}
	return strings.Join(varParts, "")
}

func replaceAllDots(s string) string {
	return strings.Replace(s, ".", "_", -1)
}
