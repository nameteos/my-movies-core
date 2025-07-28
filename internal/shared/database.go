package shared

import (
	"context"
	"fmt"
	"log"
	"time"

	"event-driven-go/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DatabaseConfig holds configuration for both databases
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

// DatabaseConnections holds both database connections
type DatabaseConnections struct {
	PostgreSQL *gorm.DB
	MongoDB    *mongo.Database
	logger     *log.Logger
}

// NewDatabaseConnections creates connections to both databases
func NewDatabaseConnections(config DatabaseConfig, logger *log.Logger) (*DatabaseConnections, error) {
	if logger == nil {
		logger = log.Default()
	}

	// PostgreSQL with GORM
	pgDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgreSQL.Host,
		config.PostgreSQL.Port,
		config.PostgreSQL.User,
		config.PostgreSQL.Password,
		config.PostgreSQL.Database,
		config.PostgreSQL.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: gormlogger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
			gormlogger.Config{
				SlowThreshold: time.Second,     // Slow SQL threshold
				LogLevel:      gormlogger.Info, // Log level
				Colorful:      true,            // Enable color
			},
		),
	}

	pgDB, err := gorm.Open(postgres.Open(pgDSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL with GORM: %w", err)
	}

	// Test the connection
	sqlDB, err := pgDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	logger.Println("✅ Connected to PostgreSQL with GORM")

	// MongoDB connection
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	mongoDB := mongoClient.Database(config.MongoDB.Database)
	logger.Println("✅ Connected to MongoDB")

	return &DatabaseConnections{
		PostgreSQL: pgDB,
		MongoDB:    mongoDB,
		logger:     logger,
	}, nil
}

// Close closes database connections
func (dc *DatabaseConnections) Close() error {
	var errors []error

	// Close PostgreSQL
	if sqlDB, err := dc.PostgreSQL.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close PostgreSQL: %w", err))
		} else {
			dc.logger.Println("🔌 Closed PostgreSQL connection")
		}
	}

	// Close MongoDB
	if err := dc.MongoDB.Client().Disconnect(context.Background()); err != nil {
		errors = append(errors, fmt.Errorf("failed to close MongoDB: %w", err))
	} else {
		dc.logger.Println("🔌 Closed MongoDB connection")
	}

	if len(errors) > 0 {
		return fmt.Errorf("database close errors: %v", errors)
	}

	return nil
}

func GetDatabaseConfig() DatabaseConfig {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	dbConfig := cfg.GetDatabaseConfig()
	return DatabaseConfig{
		PostgreSQL: PostgreSQLConfig{
			Host:     dbConfig.PostgreSQL.Host,
			Port:     dbConfig.PostgreSQL.Port,
			User:     dbConfig.PostgreSQL.User,
			Password: dbConfig.PostgreSQL.Password,
			Database: dbConfig.PostgreSQL.Database,
			SSLMode:  dbConfig.PostgreSQL.SSLMode,
		},
		MongoDB: MongoDBConfig{
			URI:      dbConfig.MongoDB.URI,
			Database: dbConfig.MongoDB.Database,
		},
	}
}
