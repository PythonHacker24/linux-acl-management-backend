package config

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

/* app parameters */
type App struct {
	Name           string `yaml:"name,omitempty"`
	Version        string `yaml:"version,omitempty"`
	DebugMode      bool   `yaml:"debug_mode,omitempty"`
	SessionTimeout int    `yaml:"session_timeout,omitempty"`
	BasePath       string `yaml:"base_path,omitempty"`
}

/* normalization function */
func (a *App) Normalize() error {

	/* set default name to laclm */
	if a.Name == "" {
		a.Name = "laclm"
	}

	/* set default version to v1.1 */
	if a.Version == "" {
		a.Name = "v1.1"
	}

	/*
		if debug_mode is not provided, it's false
		we want production to be true
	*/

	/* set default session timeout to 24 hours */
	if a.SessionTimeout == 0 {
		a.SessionTimeout = 24
	}

	/* check if base path is specified */
	if a.BasePath == "" {
		return errors.New(heredoc.Doc(`
			Base path is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	return nil
}
