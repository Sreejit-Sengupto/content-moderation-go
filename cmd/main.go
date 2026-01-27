package main

import (
	"log"
	"net/http"

	"github.com/Sreejit-Sengupto/api/routes"
	"github.com/Sreejit-Sengupto/internal/database"
	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/Sreejit-Sengupto/internal/queue"
	workerClient "github.com/Sreejit-Sengupto/internal/queue/worker-client"
	"github.com/Sreejit-Sengupto/utils/cors"
	"github.com/Sreejit-Sengupto/utils/gemini"
	"github.com/Sreejit-Sengupto/utils/imagekit"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database.ConnectDatabase()

	// Enable uuid-ossp extension
	database.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	database.DB.AutoMigrate(&models.Content{}, &models.Audit{}, &models.ModerationResult{}, &models.ModerationEvents{})

	// init imagekit
	imagekit.InitImageKit()

	// Init gemini (uses context.Background() internally for long-lived client)
	if err := gemini.InitGemini(); err != nil {
		log.Fatalf("Failed to initialize Gemini: %v", err)
		return

	}

	workerClient.InitClient()
	defer workerClient.CloseClient()

	go queue.StartWorkerServer()

	// Init router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	routes.RegisterRoutes(r)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", cors.Middleware(r)))
}
