package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

/*
	we need config normalization as well
	config normalization fixes all the fields that are not present in config file 
	and sets it to default value
*/

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

/* loads environment variables */
func LoadEnv() {

	/* get the JWT_SECRET_KEY from environment variable */
	secret := os.Getenv("JWT_SECRET_KEY")
    if secret == "" {
        zap.L().Fatal("JWT_SECRET_KEY environment variable not set")
    }
	
	EnvConfig.JWTSecret = secret
}
