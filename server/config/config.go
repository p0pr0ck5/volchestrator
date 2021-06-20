package config

type Config struct {
	ClientTTL int
}

func DefaultConfig() *Config {
	return &Config{
		ClientTTL: 30,
	}
}
