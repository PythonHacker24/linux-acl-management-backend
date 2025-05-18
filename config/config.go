package config

import "fmt"

/* globally accessible yaml config */
var BackendConfig Config

/* complete yaml config for global usage */
type Config struct {
	AppInfo           App                 `yaml:"app,omitempty"`
	Server            Server              `yaml:"server,omitempty"`
	Database          Database            `yaml:"database,omitempty"`
	Logging           Logging             `yaml:"logging,omitempty"`
	FileSystemServers []FileSystemServers `yaml:"filesystem_servers,omitempty"`
	BackendSecurity   BackendSecurity     `yaml:"backend_security,omitempty"`
	Authentication    Authentication      `yaml:"authentication,omitempty"`
}

/* complete config normalizer function */
func (c *Config) Normalize() error {
	if err := c.AppInfo.Normalize(); err != nil {
		return fmt.Errorf("app configuration error: %w", err)
	}

	if err := c.Server.Normalize(); err != nil {
		return fmt.Errorf("server configuration error: %w", err)
	}

	if err := c.Database.Normalize(); err != nil {
		return fmt.Errorf("database configuration error: %w", err)
	}

	if err := c.Logging.Normalize(); err != nil {
		return fmt.Errorf("logging configuration error: %w", err)
	}

	for i := range c.FileSystemServers {
		if err := c.FileSystemServers[i].Normalize(); err != nil {
			return fmt.Errorf("file system server [%d] error: %w", i, err)
		}
	}

	if err := c.BackendSecurity.Normalize(); err != nil {
		return fmt.Errorf("backend security configuration error: %w", err)
	}

	if err := c.Authentication.Normalize(); err != nil {
		return fmt.Errorf("authentication configuration error: %w", err)
	}

	return nil
}
