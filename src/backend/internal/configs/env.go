package configs

import "os"

type EnvConfig struct {
	Host string
	Port string
	AlloyHost string
	AlloyPort string
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

	alloyHost, ok := os.LookupEnv("ALLOY_HOST")
	if !ok {
		alloyHost = "alloy"
	}

	alloyPort, ok := os.LookupEnv("ALLOY_PORT")
	if !ok {
		alloyPort = "4318"
	}

	return &EnvConfig{
		Host: host,
		Port: port,
		AlloyHost: alloyHost,
		AlloyPort: alloyPort,
	}
}
