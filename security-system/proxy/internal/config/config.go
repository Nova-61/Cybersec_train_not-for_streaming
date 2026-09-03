package config

type Config struct {
	ListenAddr string
	TargetURL  string
}

func New() *Config {
	return &Config{
		ListenAddr: ":8888",
		TargetURL:  "http://localhost:8080",
	}
}
