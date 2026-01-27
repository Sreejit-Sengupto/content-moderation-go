package queue

import (
	"log"
	"os"

	"github.com/Sreejit-Sengupto/internal/queue/tasks"
	"github.com/Sreejit-Sengupto/internal/queue/workers/aggregation"
	"github.com/Sreejit-Sengupto/internal/queue/workers/image"
	"github.com/Sreejit-Sengupto/internal/queue/workers/text"
	"github.com/hibiken/asynq"
)

func StartWorkerServer() {
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

	srv := asynq.NewServer(
		opt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"text":        5,
				"image":       3,
				"video":       1,
				"aggregation": 1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeTextDelivery, text.HandleTextDelivery)
	mux.HandleFunc(tasks.TypeImageDelivery, image.HandleImageDelivery)
	// mux.HandleFunc(tasks.TypeVideoDelivery, text.HandleVideoDelivery)
	mux.HandleFunc(tasks.TypeAggregationDelivery, aggregation.HandleAggregationDelivery)

	log.Println("Starting Asynq worker server...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
