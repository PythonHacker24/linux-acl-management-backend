package config

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
