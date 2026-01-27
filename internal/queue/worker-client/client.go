package workerClient

import (
	"log"
	"os"

	"github.com/hibiken/asynq"
)

var Client *asynq.Client

func InitClient() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable not set")
		return
	}

	// ParseRedisURI handles rediss:// URLs (with TLS) like Upstash
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse REDIS_URL: %v", err)
		return
	}

	Client = asynq.NewClient(opt)
	log.Println("Asynq client initialized")
}

func CloseClient() {
	if Client != nil {
		Client.Close()
	}
}
