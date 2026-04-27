package config

import "os"

type Config struct {
	Port       string
	SocketPath string
}

func Load() Config {
	return Config{
		Port:       getEnv("PORT", "9999"),
		SocketPath: os.Getenv("SOCKET_PATH"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
