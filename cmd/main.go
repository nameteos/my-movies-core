package main

import (
	"context"
	"event-driven-go/internal/config"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-driven-go/internal/domains/library"
	"event-driven-go/internal/domains/movies"
	"event-driven-go/internal/domains/rating"
	"event-driven-go/internal/domains/user"
	"event-driven-go/internal/domains/watchlist"
	"event-driven-go/internal/shared"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	config.Load()

	logger := log.New(os.Stdout, "[MOVIES-GO] ", log.LstdFlags|log.Lshortfile)

	logger.Println("Starting ...")

	dbConnections, err := shared.NewDatabaseConnections(shared.GetDatabaseConfig(), logger)
	if err != nil {
		logger.Fatalf("❌ Failed to connect to databases: %v", err)
	}
	defer func(dbConnections *shared.DatabaseConnections) {
		err := dbConnections.Close()
		if err != nil {

		}
	}(dbConnections)

	userRepo := user.NewRepository(dbConnections.PostgreSQL)
	watchlistRepo := watchlist.NewRepository(dbConnections.PostgreSQL)
	libraryRepo := library.NewRepository(dbConnections.PostgreSQL)
	ratingRepo := rating.NewRepository(dbConnections.PostgreSQL)
	movieRepo := movies.NewMongoRepository(dbConnections.MongoDB)

	if err := runMigrations(userRepo, watchlistRepo, libraryRepo, ratingRepo); err != nil {
		logger.Fatalf("❌ Migration failed: %v", err)
	}
	if err := movieRepo.CreateIndexes(context.Background()); err != nil {
		logger.Fatalf("❌️  Warning: Failed to create MongoDB indexes: %v", err)
	}

	logger.Println("✅ App setup finished.")

	eventBus := shared.NewEventBus(logger)
	setupConsumers(context.Background(), logger)
	userService := user.NewService(userRepo, eventBus)

	demonstrateGormFeatures(userService, watchlistRepo, libraryRepo, ratingRepo, movieRepo, logger)

	// Keep the application running until terminated
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Println("Shutting down application...")
}

func runMigrations(
	userRepo *user.Repository,
	watchlistRepo *watchlist.Repository,
	libraryRepo *library.Repository,
	ratingRepo *rating.Repository,
) error {
	if err := userRepo.AutoMigrate(); err != nil {
		return fmt.Errorf("user migration failed: %w", err)
	}

	if err := watchlistRepo.AutoMigrate(); err != nil {
		return fmt.Errorf("watchlist migration failed: %w", err)
	}

	if err := libraryRepo.AutoMigrate(); err != nil {
		return fmt.Errorf("library migration failed: %w", err)
	}

	if err := ratingRepo.AutoMigrate(); err != nil {
		return fmt.Errorf("rating migration failed: %w", err)
	}

	return nil
}

func setupConsumers(ctx context.Context, logger *log.Logger) {

	userHandler := user.NewHandler(logger)
	watchlistHandler := watchlist.NewHandler(logger)
	libraryHandler := library.NewHandler(logger)
	ratingHandler := rating.NewHandler(logger)

	// Map topics to handlers
	handlers := map[string]shared.EventHandler{
		user.UserRegisteredEventType:             userHandler,
		user.UserUpdatedEventType:                userHandler,
		user.UserDeletedEventType:                userHandler,
		watchlist.MovieAddedToWatchlistEventType: watchlistHandler,
		library.MovieWatchedEventType:            libraryHandler,
		rating.MovieRatedEventType:               ratingHandler,
		rating.MovieUnratedEventType:             ratingHandler,
	}

	go shared.StartConsumers(ctx, handlers)
}

func demonstrateGormFeatures(
	userService *user.Service,
	watchlistRepo *watchlist.Repository,
	libraryRepo *library.Repository,
	ratingRepo *rating.Repository,
	movieRepo *movies.MongoRepository,
	logger *log.Logger,
) {
	ctx := context.Background()

	logger.Println("\n👤 PHASE 1: User Domain Operations")
	logger.Println("==================================")

	// Create sample users
	user1, err := userService.RegisterUser(ctx, "movieew32we21111313", "f1anfe32ffwe1s11fsd@1movies.com")
	if err != nil {
		logger.Printf("❌ Failed to register user 1: %v", err)
		return
	}
	logger.Printf("✅ Registered user: %s (%s)", user1.Username, user1.Email)

	user2, err := userService.RegisterUser(ctx, "cinephile", "cinephile2@example.com")
	if err != nil {
		logger.Printf("❌ Failed to register user 2: %v", err)
		return
	}
	logger.Printf("✅ Registered user: %s (%s)", user2.Username, user2.Email)

	// List users
	users, err := userService.ListUsers(ctx, 10, 0)
	if err != nil {
		logger.Printf("❌ Failed to list users: %v", err)
	} else {
		logger.Printf("📋 Total registered users: %d", len(users))
	}

	logger.Println("\n🎬 PHASE 2: Creating Movies in MongoDB")
	logger.Println("=====================================")

	// Create sample movies
	sampleMovies := []*movies.Movie{
		{
			ID:          primitive.NewObjectID(),
			Title:       "The Shawshank Redemption",
			Genre:       []string{"Drama"},
			Year:        1994,
			Director:    []string{"Frank Darabont"},
			Description: "Two imprisoned men bond over years, finding solace and redemption through common decency.",
			Duration:    142,
			Cast: []movies.CastMember{
				{Name: "Tim Robbins", Character: "Andy Dufresne", Order: 1},
				{Name: "Morgan Freeman", Character: "Ellis Redding", Order: 2},
			},
		},
		{
			ID:          primitive.NewObjectID(),
			Title:       "The Godfather",
			Genre:       []string{"Crime", "Drama"},
			Year:        1972,
			Director:    []string{"Francis Ford Coppola"},
			Description: "The aging patriarch of a crime dynasty transfers control to his reluctant son.",
			Duration:    175,
			Cast: []movies.CastMember{
				{Name: "Marlon Brando", Character: "Don Vito Corleone", Order: 1},
				{Name: "Al Pacino", Character: "Michael Corleone", Order: 2},
			},
		},
	}

	// Store movies in MongoDB
	var createdMovies []*movies.Movie
	for _, movie := range sampleMovies {
		created, err := movieRepo.CreateMovie(ctx, movie)
		if err != nil {
			logger.Printf("❌ Failed to create movie %s: %v", movie.Title, err)
			continue
		}
		createdMovies = append(createdMovies, created)
		logger.Printf("📽️  Created: %s (%d)", created.Title, created.Year)
	}

	logger.Println("\n📋 PHASE 3: GORM Watchlist Operations")
	logger.Println("===================================")

	// Add movies to user's watchlist
	for i, movie := range createdMovies {
		notes := fmt.Sprintf("Must watch #%d - heard amazing things!", i+1)
		_, err := watchlistRepo.AddToWatchlist(ctx, user1.ID, movie.IDString(), notes)
		if err != nil {
			logger.Printf("❌ Failed to add %s to watchlist: %v", movie.Title, err)
			continue
		}
		logger.Printf("✅ Added to %s's watchlist: %s", user1.Username, movie.Title)
	}

	// Get user's watchlist
	userWatchlist, err := watchlistRepo.GetUserWatchlist(ctx, user1.ID)
	if err != nil {
		logger.Printf("❌ Failed to get watchlist: %v", err)
	} else {
		logger.Printf("📋 %s has %d movies in watchlist", user1.Username, len(userWatchlist))
	}

	logger.Println("\n📚 PHASE 4: GORM Library Operations")
	logger.Println("==================================")

	// Mark movies as watched
	for i, movie := range createdMovies {
		watchedAt := time.Now().Add(-time.Duration(i+1) * 24 * time.Hour)
		_, err := libraryRepo.AddWatchHistory(ctx, user1.ID, movie.IDString(), watchedAt, movie.Duration)
		if err != nil {
			logger.Printf("❌ Failed to add watch history for %s: %v", movie.Title, err)
			continue
		}
		logger.Printf("📺 %s watched: %s", user1.Username, movie.Title)
	}

	// Get watching statistics
	stats, err := libraryRepo.GetWatchingStats(ctx, user1.ID)
	if err != nil {
		logger.Printf("❌ Failed to get stats: %v", err)
	} else {
		logger.Printf("📊 %s's Stats: Total movies: %v, Total hours: %.1f",
			user1.Username, stats["total_movies"], stats["total_hours"])
	}

	logger.Println("\n⭐ PHASE 5: GORM Rating Operations")
	logger.Println("=================================")

	// Rate the watched movies
	ratings := []struct {
		movie  *movies.Movie
		rating float64
		review string
	}{
		{createdMovies[0], 4.8, "Absolutely incredible! One of the best films ever made."},
		{createdMovies[1], 4.9, "A masterpiece of cinema. Perfect storytelling and acting."},
	}

	for _, r := range ratings {
		movieRating, err := ratingRepo.UpsertRating(ctx, user1.ID, r.movie.IDString(), r.rating, r.review)
		if err != nil {
			logger.Printf("❌ Failed to rate %s: %v", r.movie.Title, err)
			continue
		}
		logger.Printf("⭐ %s rated: %s - %.1f/5", user1.Username, r.movie.Title, movieRating.Rating)
	}

	// Update user profile
	user1.Username = "moviefan123_updated"
	updatedUser, err := userService.UpdateUser(ctx, user1)
	if err != nil {
		logger.Printf("❌ Failed to update user: %v", err)
	} else {
		logger.Printf("✅ Updated user profile: %s", updatedUser.Username)
	}

	logger.Println("\n🎉 GORM Demo with User Domain Completed Successfully!")
	logger.Println("===================================================")
	logger.Printf("🔍 Key Features Demonstrated:")
	logger.Printf("✅ User domain with complete CRUD operations")
	logger.Printf("✅ Event-driven user lifecycle management")
	logger.Printf("✅ Cross-domain data relationships")
	logger.Printf("✅ GORM auto-migrations and type safety")
	logger.Printf("✅ MongoDB + PostgreSQL hybrid architecture")
	logger.Printf("✅ Domain-driven design principles")
}
