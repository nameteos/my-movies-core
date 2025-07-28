package config

import (
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	PostgresHost     string `env:"POSTGRES_HOST" default:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT" default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" default:"movieapp"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" default:"movieapp123"`
	PostgresDB       string `env:"POSTGRES_DB" default:"movieapp"`
	PostgresSSLMode  string `env:"POSTGRES_SSLMODE" default:"disable"`

	MongoURI      string `env:"MONGODB_URI" default:"mongodb://localhost:27017"`
	MongoDatabase string `env:"MONGODB_DATABASE" default:"movieapp"`

	Environment string `env:"APP_ENV" default:"development"`
	LogLevel    string `env:"LOG_LEVEL" default:"info"`

	KafkaBootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS" default:"localhost:9092"`
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Try to load .env file, but continue if it doesn't exist
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables and defaults")
	}

	config := &Config{}
	if err := parseEnv(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}

func parseEnv(config interface{}) error {
	configValue := reflect.ValueOf(config)
	if configValue.Kind() != reflect.Ptr {
		return errors.New("config must be a pointer to a struct")
	}

	configElem := configValue.Elem()
	if configElem.Kind() != reflect.Struct {
		return errors.New("config must be a pointer to a struct")
	}

	configType := configElem.Type()
	for i := 0; i < configElem.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := configElem.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		envName := field.Tag.Get("env")
		if envName == "" {
			continue
		}

		defaultValue := field.Tag.Get("default")
		required := field.Tag.Get("required") == "true"

		envValue := os.Getenv(envName)
		if envValue == "" {
			if required {
				return fmt.Errorf("required environment variable %s is not set", envName)
			}
			envValue = defaultValue
		}

		if err := setField(fieldValue, envValue); err != nil {
			return fmt.Errorf("failed to set field %s: %w", field.Name, err)
		}
	}

	return nil
}

// setField sets the value of a struct field based on the environment variable value
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			value = "0"
		}
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int value: %w", err)
		}
		field.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			value = "0"
		}
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint value: %w", err)
		}
		field.SetUint(uintValue)
	case reflect.Bool:
		if value == "" {
			value = "false"
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to parse bool value: %w", err)
		}
		field.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			value = "0"
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float value: %w", err)
		}
		field.SetFloat(floatValue)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			if value == "" {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			} else {
				values := strings.Split(value, ",")
				slice := reflect.MakeSlice(field.Type(), len(values), len(values))
				for i, v := range values {
					slice.Index(i).SetString(strings.TrimSpace(v))
				}
				field.Set(slice)
			}
		} else {
			return fmt.Errorf("unsupported slice type: %s", field.Type().Elem().Kind())
		}
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if value == "" {
				field.Set(reflect.ValueOf(time.Time{}))
			} else {
				timeValue, err := time.Parse(time.RFC3339, value)
				if err != nil {
					return fmt.Errorf("failed to parse time value: %w", err)
				}
				field.Set(reflect.ValueOf(timeValue))
			}
		} else {
			return fmt.Errorf("unsupported struct type: %s", field.Type().Name())
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		PostgreSQL: PostgreSQLConfig{
			Host:     c.PostgresHost,
			Port:     c.PostgresPort,
			User:     c.PostgresUser,
			Password: c.PostgresPassword,
			Database: c.PostgresDB,
			SSLMode:  c.PostgresSSLMode,
		},
		MongoDB: MongoDBConfig{
			URI:      c.MongoURI,
			Database: c.MongoDatabase,
		},
	}
}

func (c *Config) GetSeranaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Consumer.Return.Errors = true

	return config
}

type DatabaseConfig struct {
	PostgreSQL PostgreSQLConfig
	MongoDB    MongoDBConfig
}

type PostgreSQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type MongoDBConfig struct {
	URI      string
	Database string
}
