package config

/* server deployment parameters */
type Server struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

/* normalization function */
func (s *Server) Normalize() error {

	/* set default host to localhost */
	if s.Host == "" {
		s.Host = "localhost"
	}

	/* set default port to 8080 */
	if s.Port == 0 {
		s.Port = 8080
	}

	return nil
}
