package config

/* app parameters */
type App struct {
Name        string `yaml:"name"`
Version     string `yaml:"version"`
Environment string `yaml:"environment"`
}
