package common

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Database *DatabaseConfig
}

type DatabaseConfig struct {
	Driver string

	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (dc *DatabaseConfig) URI() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s",
		dc.Driver,
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.Database,
	)
}

func ConfigFromEnv() *Config {
	dbPort, err := strconv.Atoi(getEnv("POSTGRES_PORT"))
	if err != nil {
		panic("invalid port variable type")
	}
	return &Config{
		Database: &DatabaseConfig{
			Driver:   "postgresql",
			Host:     getEnv("POSTGRES_HOST"),
			Port:     dbPort,
			User:     getEnv("POSTGRES_USER"),
			Password: getEnv("POSTGRES_PASSWORD"),
			Database: getEnv("POSTGRES_DB"),
		},
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("The %q environment variable not set", key))
	}

	return value
}
