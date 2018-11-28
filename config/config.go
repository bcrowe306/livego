package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Data :This variable holds the application configuration data.
var Data *Config

// Config :This struct represents the Application configuration
type Config struct {
	RTMPPort string            `yaml:"rtmpport"`
	APIPort  string            `yaml:"apiport"`
	Backend  map[string]string `yaml:"backend"`
}

// LoadConfig :This function loads the configuration from the specified file.
func LoadConfig(fileName string) error {
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(fileBytes, &Data); err != nil {
		return err
	}
	return nil
}

// SaveConfig :This function saves the application configuration to the file specified.
func SaveConfig(fileName string) error {
	fileBytes, err := yaml.Marshal(Data)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, fileBytes, 0777); err != nil {
		return err
	}
	return nil
}
