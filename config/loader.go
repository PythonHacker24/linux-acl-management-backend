package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"github.com/davecgh/go-spew/spew"
)

/*
	we need config normalization as well
	config normalization fixes all the fields that are not present in config file
	and sets it to default value
*/

/* loads yaml config file from given file path */
func LoadConfig(path string) error {

	/* read the yaml config file */
    data, err := os.ReadFile(path)
    if err != nil {
		return fmt.Errorf("config loading error %w", 
			err,
		)

    }

	/* unmarshal the yaml file to defined struct */
    err = yaml.Unmarshal(data, &BackendConfig)
    if err != nil {
		return fmt.Errorf("config loading error %w", 
			err,
		)
    }

	if BackendConfig.AppInfo.DebugMode {
		fmt.Println("Contents of Config File (debug mode ON)")
		spew.Dump(BackendConfig)
		fmt.Println()
	}
	
	/* normalize the complete backend config before proceeding */
	return BackendConfig.Normalize()
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
