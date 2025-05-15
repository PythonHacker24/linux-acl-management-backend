package config

/* file system server parameters */
type FileSystemServers struct {
	Remote *Remote `yaml:"remote,omitempty"`
	Path   string  `yaml:"path"`
	Method string  `yaml:"method"`
}

/* remote parameters for file system server with laclm daemons installed */
type Remote struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
