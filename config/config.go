package config

/* globally accessible yaml config */
var BackendConfig Config

/* globally accessible environment variables */
var EnvConfig EnvironmentConfig

/* complete yaml config for global usage */
type Config struct {
	AppInfo           App                 `yaml:"app"`
	Server            Server              `yaml:"server"`
	Database          Database            `yaml:"database"`
	Logging           Logging             `yaml:"logging"`
	FileSystemServers []FileSystemServers `yaml:"filesystem_servers"`
	BackendSecurity   BackendSecurity     `yaml:"backend_security"`
}

/* complete environment variables configs for global usage */
type EnvironmentConfig struct {
	JWTSecret string
}
