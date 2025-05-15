package config

/* backend security configs */
type BackendSecurity struct {
	JWTExpiry int `yaml:"jwt_expiry"`
}

/* normalization function */
func (b *BackendSecurity) Normalize() error {
	if b.JWTExpiry == 0 {
		b.JWTExpiry = 1
	}

	return nil
}
