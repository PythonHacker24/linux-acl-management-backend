package config

import (
	"fmt"
	"os"

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

	/* expand all environment variables in the yaml config */
	expanded := os.ExpandEnv(string(data))

	/* unmarshal the yaml file to defined struct */
    err = yaml.Unmarshal([]byte(expanded), &BackendConfig)
    if err != nil {
		return fmt.Errorf("config loading error %w", 
			err,
		)
    }

	/* write the config file in console if in debug mode */
	if BackendConfig.AppInfo.DebugMode {
		fmt.Println("Contents of Config File (debug mode ON)")
		spew.Dump(BackendConfig)
		fmt.Println()
	}
	
	/* normalize the complete backend config before proceeding */
	return BackendConfig.Normalize()
}
