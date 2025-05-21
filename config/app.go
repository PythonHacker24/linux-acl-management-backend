package config

/* app parameters */
type App struct {
	Name      string `yaml:"name,omitempty"`
	Version   string `yaml:"version,omitempty"`
	DebugMode bool   `yaml:"debug_mode,omitempty"`
	SessionTimeout	int	`yaml:"session_timeout,omitempty"`
}

/* normalization function */
func (a *App) Normalize() error {
	if a.Name == "" {
		a.Name = "laclm"
	}

	if a.Version == "" {
		a.Name = "v1.1"
	}

	/*
		if debug_mode is not provided, it's false
		we want production to be true
	*/

	if a.SessionTimeout == 0 {
		a.SessionTimeout = 24
	}

	return nil
}
