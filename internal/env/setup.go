package env

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

type postgresConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type Config struct {
	Port        int
	Debug       bool
	ExternalAPI string
	Postgres    postgresConfig
}

func Setup() *Config {
	log.Debugln("initialize config")
	return &Config{
		Port:        getEnvAsInt("PORT", 8080),
		Debug:       getEnvAsBool("DEBUG_MODE", true),
		ExternalAPI: getEnv("EXTERNAL_API", ""),
		Postgres: postgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_HOST_PORT", 5432),
			Database: getEnv("POSTGRES_DB", "default"),
			User:     getEnv("POSTGRES_USER", "root"),
			Password: getEnv("POSTGRES_PASSWORD", "root"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}
