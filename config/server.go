package config

/* server deployment parameters */
type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

/* normalization function */
func (s *Server) Normalize() error {
	if s.Host == "" {
		s.Host = "localhost"
	}

	if s.Port == 0 {
		s.Port = 8080
	}

	return nil
}
