package config

/* backend security configs */
type BackendSecurity struct {
	JWTExpiry int `yaml:"jwt_expiry"`
}
