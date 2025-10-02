package config

import (
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	LoggerConfigKey     = "LOG"
	MetricsConfigKey    = "METRICS"
	CalculatorConfigKey = "CALCULATOR"
)

type Configs map[string]interface{}

type CalculatorConfig struct {
	Version         string        `json:"version"`
	StorageType     string        `json:"storage_type"`
	StorageFilePath string        `json:"storage_file_path"`
	Port            string        `json:"port"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
}
type LoggerConfig struct {
	ServerName string `json:"server_name"`
	TimeFormat string `json:"time_format"`
	Version    string `json:"version"`
	Level      string `json:"level"`
	Format     string `json:"format"` // json or text
}

type MetricsConfig struct {
	ServiceName    string `json:"service_name"`
	ServiceVersion string `json:"service_version"`
	ServerName     string `json:"server_name"`
}

// LoadConfigs is a pseudo factory pattern to load requested configurations
// a bit overkill for this case, but could be useful in a more complex scenario
func LoadConfigs(requestedServices []string) Configs {
	configs := make(map[string]interface{})
	serverName := getEnv("SERVER_NAME", "unknown_server")
	for _, service := range requestedServices {
		switch service {
		case LoggerConfigKey:
			logConfig := getLoggerConfig()
			logConfig.ServerName = serverName
			configs[LoggerConfigKey] = logConfig
		case MetricsConfigKey:
			metricsConfig := MetricsConfig{
				ServerName: serverName,
			}
			configs[MetricsConfigKey] = metricsConfig
		case CalculatorConfigKey:
			calculatorConfig := getDefaultCalculatorConfig()
			configs[CalculatorConfigKey] = calculatorConfig
		}
	}
	return configs
}

func getLoggerConfig() LoggerConfig {
	timeFormat := getEnv("LOG_TIME_FORMAT", "2006-01-02 15:04:05")
	level := getEnv("LOG_LEVEL", logrus.InfoLevel.String())
	format := getEnv("LOG_FORMAT", "text")

	return LoggerConfig{
		TimeFormat: timeFormat,
		Level:      level,
		Format:     format,
	}
}
func getDefaultCalculatorConfig() CalculatorConfig {
	version := getEnv("CALCULATOR_VERSION", "1.0.0")
	port := getEnv("CALCULATOR_PORT", "8080")
	storageType := getEnv("CALCULATOR_STORAGE_TYPE", "memory")
	storageFilePath := getEnv("CALCULATOR_STORAGE_PATH", "./storage.txt")
	readTimeout := time.Second * time.Duration(getEnvAsInt("CALCULATOR_READ_TIMEOUT", 5))
	writeTimeout := time.Second * time.Duration(getEnvAsInt("CALCULATOR_WRITE_TIMEOUT", 10))
	idleTimeout := time.Second * time.Duration(getEnvAsInt("CALCULATOR_IDLE_TIMEOUT", 120))

	return CalculatorConfig{
		Version:         version,
		StorageType:     storageType,
		StorageFilePath: storageFilePath,
		Port:            port,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		IdleTimeout:     idleTimeout,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
