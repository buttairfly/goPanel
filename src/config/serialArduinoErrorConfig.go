package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type ArduinoErrorConfig map[string]ArduinoErrorDescription

type ArduinoErrorDescription struct {
	Name      string `json:"name"`
	Param     string `json:"param,omitempty"`
	Character string `json:"character,omitempty"`
}

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

func (aec *ArduinoErrorConfig) GetDescription(errorCode string) (*ArduinoErrorDescription, error) {
	description, ok := (*aec)[errorCode]
	if !ok {
		return nil, fmt.Errorf("Unkown error code %s", errorCode)
	}
	return &description, nil
}

func (aec *ArduinoErrorConfig) ToCppFile() []byte {
	return nil
}
