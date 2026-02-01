package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	if os.Getenv("RUN_MIGRATION") == "TRUE" {
		database.DB.AutoMigrate(&models.Content{}, &models.Audit{}, &models.ModerationResult{}, &models.ModerationEvents{})
	}

	// init imagekit
	imagekit.InitImageKit()

	// Init gemini (uses context.Background() internally for long-lived client)
	if err := gemini.InitGemini(); err != nil {
		log.Fatalf("Failed to initialize Gemini: %v", err)
		return
	}

	workerClient.InitClient()
	defer workerClient.CloseClient()

	// Channel to signal worker server shutdown
	workerShutdown := make(chan struct{})

	// WaitGroup to track goroutines
	var wg sync.WaitGroup

	// Start worker server
	wg.Add(1)
	go func() {
		defer wg.Done()
		queue.StartWorkerServer(workerShutdown)
	}()

	// Init router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	routes.RegisterRoutes(r)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: cors.Middleware(r),
	}

	// Start HTTP server in goroutine
	go func() {
		log.Println("Server started on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-quit
	log.Printf("Received signal %v, initiating graceful shutdown...", sig)

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	log.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Signal worker server to shutdown
	log.Println("Shutting down worker server...")
	close(workerShutdown)

	// Wait for worker to finish with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All workers shut down gracefully")
	case <-ctx.Done():
		log.Println("Shutdown timed out")
	}

	log.Println("Server exited")
}
