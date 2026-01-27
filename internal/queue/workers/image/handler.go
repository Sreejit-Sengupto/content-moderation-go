package image

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Sreejit-Sengupto/internal/database"
	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/Sreejit-Sengupto/internal/queue/tasks"
	workerClient "github.com/Sreejit-Sengupto/internal/queue/worker-client"
	"github.com/Sreejit-Sengupto/utils/gemini"
	"github.com/hibiken/asynq"
	"google.golang.org/genai"
)

type ImageModerationResult struct {
	Status      string  `json:"status"`
	RiskScore   float64 `json:"riskScore"`
	Explanation string  `json:"explanation"`
}

type eventPayload struct {
	ImageURL string `json:"imageURL"`
}

const SystemInstruction = " You are an automated image content moderation system designed to evaluate user-submitted images for safety and policy compliance. Your role is to objectively assess the visual content of each image and determine whether it is suitable for publication on a public platform. You must analyze images for the presence of unsafe or prohibited visual material, including but not limited to violence, graphic injury, sexual or pornographic content, child exploitation, hate symbols, harassment, self-harm, illegal activities, extremist imagery, misleading or manipulated media, and other harmful or policy-violating elements. Your evaluation must be based only on what is visible in the image itself, without assuming intent, narrative context, or external metadata unless explicitly provided as part of the image. If an image clearly violates safety standards, it must be rejected. If an image is ambiguous, borderline, or context-dependent, it must be flagged for human review. If an image does not present any safety or policy concerns, it must be approved. Your decisions must be consistent, conservative, and explainable. Do not modify, enhance, censor, describe creatively, or interpret the image beyond safety evaluation. Do not provide advice, opinions, captions, or alternative representations. Your task is strictly limited to classification and moderation decision-making."

func HandleImageDelivery(ctx context.Context, t *asynq.Task) error {
	fmt.Println("Processing image moderation task")
	var payload tasks.ImageDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// JSON Schema for structured output
	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"status": {
				Type: genai.TypeString,
				Enum: []string{"APPROVED", "REJECTED", "FLAGGED"},
			},
			"riskScore": {
				Type:        genai.TypeNumber,
				Description: "A score between 0 and 1 indicating the risk level",
			},
			"explanation": {
				Type:        genai.TypeString,
				Description: "A brief explanation of the moderation decision",
			},
		},
		Required: []string{"status", "riskScore", "explanation"},
	}

	config := &genai.GenerateContentConfig{
		Temperature:       genai.Ptr(float32(0)),
		ResponseMIMEType:  "application/json",
		ResponseSchema:    schema,
		SystemInstruction: genai.NewContentFromText(SystemInstruction, genai.RoleUser),
	}

	// fetch image
	imageRes, err := http.Get(payload.Image)
	if err != nil {
		return fmt.Errorf("http.Get failed: %v: %w", err, asynq.SkipRetry)
	}
	defer imageRes.Body.Close()

	imageBytes, err := io.ReadAll(imageRes.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll failed: %v: %w", err, asynq.SkipRetry)
	}

	parts := []*genai.Part{
		genai.NewPartFromBytes(imageBytes, "image/jpeg"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	response, err := gemini.GeminiClient.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		contents,
		config,
	)
	if err != nil {
		return fmt.Errorf("gemini.GeminiClient.Models.GenerateContent failed: %v: %w", err, asynq.SkipRetry)
	}

	var result ImageModerationResult
	if err := json.Unmarshal([]byte(response.Text()), &result); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	db := database.DB

	moderationResult := models.ModerationResult{
		ContentId:    payload.ContentID,
		MediaType:    models.MediaType(models.Img),
		Status:       models.ContentStatus(result.Status),
		RiskScore:    result.RiskScore,
		Explaination: result.Explanation,
	}
	db.Create(&moderationResult)

	modDataPayload := eventPayload{
		ImageURL: payload.Image,
	}
	modDataEventJson, err := json.Marshal(modDataPayload)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v: %w", err, asynq.SkipRetry)
	}

	moderationEventData := models.ModerationEvents{
		ContentId: payload.ContentID,
		EventType: models.EventType(models.Moderated),
		Payload:   modDataEventJson,
	}
	db.Create(&moderationEventData)

	status := models.ContentStatus(result.Status)
	task, err := tasks.NewAggregationDeliveryTask(payload.ContentID, nil, &status, nil)
	if err != nil {
		return fmt.Errorf("tasks.NewAggregationDeliveryTask failed: %v: %w", err, asynq.SkipRetry)
	}
	info, err := workerClient.Client.Enqueue(task, asynq.Queue(tasks.QueueAggregation))
	if err != nil {
		return fmt.Errorf("workerClient.Client.Enqueue failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	return nil
}
