package config

/* server deployment parameters */
type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
