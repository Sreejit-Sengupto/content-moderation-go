package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Sreejit-Sengupto/internal/database"
	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/Sreejit-Sengupto/internal/queue/tasks"
	workerClient "github.com/Sreejit-Sengupto/internal/queue/worker-client"
	"github.com/Sreejit-Sengupto/utils/imagekit"
	"github.com/Sreejit-Sengupto/utils/response"
	"github.com/Sreejit-Sengupto/utils/validator"
	"github.com/hibiken/asynq"
)

func UploadContent(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Text  string `json:"text"`
		Image string `json:"image"`
		Video string `json:"video"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validtor().Struct(reqBody); err != nil {
		response.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	if reqBody.Text != "" {
		response.JSONError(w, http.StatusBadRequest, "Text is required")
	}

	db := database.DB

	newContent := models.Content{
		Text:  reqBody.Text,
		Image: reqBody.Image,
		Video: reqBody.Video,
	}

	result := db.Create(&newContent)
	if result.Error != nil {
		response.JSONError(w, http.StatusInternalServerError, "Failed to create content")
		return
	}

	// send required to queue
	fmt.Println("Pushing text moderation task to queue")
	task, err := tasks.NewTextDeliveryTask(newContent.ID, newContent.Text)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "Failed to create text moderation task")
		return
	}
	info, err := workerClient.Client.Enqueue(task, asynq.Queue(tasks.QueueText))
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "Failed to enqueue text moderation task")
		return
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	if newContent.Image != "" {
		fmt.Println("Pushing image moderation task to queue")
		task, err := tasks.NewImageDeliveryTask(newContent.ID, newContent.Image)
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, "Failed to create image moderation task")
			return
		}
		info, err := workerClient.Client.Enqueue(task, asynq.Queue(tasks.QueueImage))
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, "Failed to enqueue image moderation task")
			return
		}
		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

	if newContent.Video != "" {
		fmt.Println("Pushing video moderation task to queue")
		task, err := tasks.NewVideoDeliveryTask(newContent.ID, newContent.Video)
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, "Failed to create video moderation task")
			return
		}
		info, err := workerClient.Client.Enqueue(task, asynq.Queue(tasks.QueueVideo))
		if err != nil {
			response.JSONError(w, http.StatusInternalServerError, "Failed to enqueue video moderation task")
			return
		}
		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	}

	response.JSON(w, http.StatusCreated, newContent)
}

func GetImageKitParams(w http.ResponseWriter, r *http.Request) {
	authParams, err := imagekit.ImageKitClient.Helper.GetAuthenticationParameters("", 0)
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "Failed to get authentication parameters")
		return
	}
	response.JSON(w, http.StatusOK, authParams)
}
