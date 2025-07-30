package shared

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DatabaseConnections struct {
	PostgreSQL *gorm.DB
	MongoDB    *mongo.Database
	logger     *log.Logger
}

func NewDatabaseConnections() (*DatabaseConnections, error) {
	pgDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		Config.Postgres.Host,
		Config.Postgres.Port,
		Config.Postgres.User,
		Config.Postgres.Password,
		Config.Postgres.Database,
		Config.Postgres.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: gormlogger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags),
			gormlogger.Config{
				SlowThreshold: time.Second,
				LogLevel:      gormlogger.Info,
				Colorful:      true,
			},
		),
	}

	pgDB, err := gorm.Open(postgres.Open(pgDSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL with GORM: %w", err)
	}

	sqlDB, err := pgDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL with GORM")

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(Config.MongoDB.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	mongoDB := mongoClient.Database(Config.MongoDB.Database)
	log.Println("✅ Connected to MongoDB")

	return &DatabaseConnections{
		PostgreSQL: pgDB,
		MongoDB:    mongoDB,
	}, nil
}

func (dc *DatabaseConnections) Close() error {
	var errors []error

	if sqlDB, err := dc.PostgreSQL.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close PostgreSQL: %w", err))
		} else {
			dc.logger.Println("🔌 Closed PostgreSQL connection")
		}
	}

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
