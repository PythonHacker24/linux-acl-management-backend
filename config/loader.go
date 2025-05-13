package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

/* loads yaml config file from given file path */
func LoadConfig(path string) {

	/* read the yaml config file */
    data, err := os.ReadFile(path)
    if err != nil {
        zap.L().Fatal("Failed to read config file",
			zap.Error(err),
		)
    }

	/* unmarshal the yaml file to defined struct */
    err = yaml.Unmarshal(data, &BackendConfig)
    if err != nil {
        zap.L().Fatal("Failed to parse YAML config", 
			zap.Error(err),
		)
    }
}
