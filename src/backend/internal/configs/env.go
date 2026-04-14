package configs

import "os"

type EnvConfig struct {
	Host string
	Port string
}

func NewEnvConfig() *EnvConfig {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "localhost"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	return &EnvConfig{
		Host: host,
		Port: port,
	}
}
