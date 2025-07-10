package vars

import (
	"fmt"
	"github.com/projectdiscovery/gologger"
	"gopkg.in/yaml.v3"
	"os"
)

type P1fingerConf struct {
	RuleMode string `yaml:"RuleMode"`

	FofaCredentials struct {
		Email  string `yaml:"Email"`
		ApiKey string `yaml:"ApiKey"`
	} `yaml:"FofaCredentials"`
}

func LoadAppConf(filePath string, config *P1fingerConf) error {

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		defaultConfig := P1fingerConf{
			RuleMode: "redteam",
			FofaCredentials: struct {
				Email  string `yaml:"Email"`
				ApiKey string `yaml:"ApiKey"`
			}{
				Email:  "P001water@163.com",
				ApiKey: "xxxx",
			},
		}
		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			return fmt.Errorf("An error occurred when generating the default configuration: %v", err)
		}

		err = os.WriteFile(filePath, data, 0644)
		if err != nil {
			return fmt.Errorf("It is impossible to create a file and write the default configuration: %v", err)
		}

		gologger.Info().Msgf("The configuration file does not exist. A file has been created in the current directory and the default configuration has been written")
		gologger.Info().Msgf("File path: %s", filePath)
		os.Exit(0)
	} else if err != nil {
		return fmt.Errorf("An error occurred when checking the file status: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("The configuration file cannot be read: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("An error occurred when parsing the configuration file: %v", err)
	}

	return nil
}
