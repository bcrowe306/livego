package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Data :This variable holds the application configuration data.
var Data *Config
var configFile *string

// Config :This struct represents the Application configuration
type Config struct {
	RTMPPort string            `yaml:"rtmpport"`
	APIPort  string            `yaml:"apiport"`
	HLSPort  string            `yaml:"hlsport"`
	Backend  map[string]string `yaml:"backend"`
	Username string            `yaml:"username"`
	Password string            `yaml:"password"`
}

var defaultConfig = Config{
	RTMPPort: "1935",
	APIPort:  "8080",
	HLSPort:  "7002",
	Backend: map[string]string{
		"name": "yaml",
		"file": "./routes.yaml",
	},
	Username: "admin",
	Password: "password",
}

// LoadConfig :This function loads the configuration from the specified file.
func LoadConfig(fileName string) error {
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		// If the config file does not exist we create it with default options.
		if os.IsNotExist(err) {
			log.Println("Creating configuration file with default values")
			outBytes, err := yaml.Marshal(defaultConfig)
			if err != nil {
				log.Println("Error Marshalling defaultConfig to outBytes")
				return err
			}
			if err := ioutil.WriteFile(fileName, outBytes, 0777); err != nil {
				log.Panicln("Error writing default config file with default options.")
				return err
			}
			Data = &defaultConfig
		}
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
