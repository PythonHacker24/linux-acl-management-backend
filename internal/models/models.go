package models

/* app parameters */
type App struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

/* server deployment parameters */
type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

/* database parameters */
type Database struct {
	TransactionLogRedis TransactionLogRedis `yaml:"transaction_log_redis"`
}

/* transaction log redis parameters */
type TransactionLogRedis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}

/* logging parameters */
type Logging struct {
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   int    `yaml:"compress"`
}

/* file system server parameters */
type FileSystemServers struct {
	Path   string  `yaml:"path"`
	Method string  `yaml:"method"`
	Remote *Remote `yaml:"remote"`
}

/* remote parameters for file system server with laclm daemons installed */
type Remote struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

/* backend security configs */
type BackendSecurity struct {
	JWTExpiry int `yaml:"jwt_expiry"`
}

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
type EnvConfig struct {
	JWTSecret string
}

/* health response */
type HealthResponse struct {
	Status string `json:"status"`
}

/* username and password */
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
