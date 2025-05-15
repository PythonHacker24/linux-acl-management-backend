package models

func (m *Config) Normalize() {
	if m.Server.Host == "" {
		m.Server.Host = "localhost" 
	}

	if m.Server.Port == 0 {
		m.Server.Port = 8080
	}
}
