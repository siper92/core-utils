package core_utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

func LoadStructFromFile(s interface{}, yamlPath string) error {
	yamlFileContent, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}

	return LoadStructFromBytes(s, yamlFileContent)
}

func LoadStructFromBytes(s interface{}, data []byte) error {
	var err error

	if err = yaml.Unmarshal(data, s); err != nil {
		return err
	}

	return ValidateStruct(s)
}
