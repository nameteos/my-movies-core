package shared

import (
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var Config *config

type config struct {
	Postgres PostgreSQLConfig
	MongoDB  MongoDBConfig
	Kafka    KafkaConfig
	App      AppConfig
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

type KafkaConfig struct {
	BootstrapServers string
	ConsumerGroup    string
	Enabled          bool
}

type AppConfig struct {
	Environment string
	LogLevel    string
}

func init() {
	Config = newConfig()
}

func newConfig() *config {
	conf, err := load()
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

func load() (*config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file, using environment variables and defaults")
	}

	config := &config{
		Postgres: PostgreSQLConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			User:     getEnv("POSTGRES_USER", "movieapp"),
			Password: getEnv("POSTGRES_PASSWORD", "movieapp123"),
			Database: getEnv("POSTGRES_DB", "movieapp"),
			SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "movieapp"),
		},
		Kafka: KafkaConfig{
			BootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
			ConsumerGroup:    getEnv("KAFKA_CONSUMER_GROUP", "movieapp"),
			Enabled:          getEnvAsBool("KAFKA_ENABLED", true),
		},
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	values := strings.Split(valueStr, sep)
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return values
}

func (c *config) GetSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Consumer.Return.Errors = true

	return config
}
