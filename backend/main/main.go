package main

import (
	"context"
	"log"
	"os"
	"time"

	"bug-tracker/controller"
	"bug-tracker/repository"
	"bug-tracker/router"
	"bug-tracker/usecase"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get environment variables
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "bug_tracker")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create client options with proper settings for Atlas
	clientOptions := options.Client().ApplyURI(mongoURI).
		SetServerSelectionTimeout(5 * time.Second).
		SetConnectTimeout(10 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Successfully connected to MongoDB!")

	// Get database
	db := client.Database(dbName)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	bugRepo := repository.NewBugRepository(db)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtSecret)
	bugUseCase := usecase.NewBugUseCase(bugRepo, userRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authUseCase)
	bugController := controller.NewBugController(bugUseCase)

	// Initialize router
	r := router.NewRouter(authController, bugController, authUseCase)
	router := r.Setup()

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
