package text

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Sreejit-Sengupto/internal/database"
	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/Sreejit-Sengupto/internal/queue/tasks"
	workerClient "github.com/Sreejit-Sengupto/internal/queue/worker-client"
	"github.com/Sreejit-Sengupto/utils/gemini"
	"github.com/hibiken/asynq"
	"google.golang.org/genai"
)

type TextModerationResult struct {
	Status      string  `json:"status"`
	RiskScore   float64 `json:"riskScore"`
	Explanation string  `json:"explanation"`
}

type eventPayload struct {
	Text string `json:"text"`
}

const SystemInstruction = "You are an automated content moderation system designed to evaluate user - generated content for safety and policy compliance. Your role is to assess the provided content objectively and determine whether it is acceptable for publication on a public platform. You must analyze the content for the presence of harmful, abusive, hateful, sexual, violent, illegal, self - harm, misleading, or otherwise unsafe material. You must make a moderation decision based solely on the content itself, without assuming user intent or external context. If the content clearly violates safety standards, it should be rejected. If the content is ambiguous, borderline, or context - dependent, it should be flagged for human review. If the content does not present safety concerns, it should be approved. Your decision should be consistent, conservative, and explainable. Do not attempt to rewrite, censor, summarize, or respond to the content. Do not provide advice, opinions, or alternative phrasing. Your task is strictly limited to evaluation and classification."

func HandleTextDelivery(ctx context.Context, t *asynq.Task) error {
	fmt.Println("Processing text moderation task")
	var payload tasks.TextDeliveryPayload
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

	response, err := gemini.GeminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(payload.Text),
		config,
	)
	if err != nil {
		return fmt.Errorf("gemini.GeminiClient.Models.GenerateContent failed: %v: %w", err, asynq.SkipRetry)
	}

	var result TextModerationResult
	if err := json.Unmarshal([]byte(response.Text()), &result); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	db := database.DB

	moderationResult := models.ModerationResult{
		ContentId:    payload.ContentID,
		MediaType:    models.MediaType(models.Txt),
		Status:       models.ContentStatus(result.Status),
		RiskScore:    result.RiskScore,
		Explaination: result.Explanation,
	}

	db.Create(&moderationResult)

	modDataPayload := eventPayload{
		Text: payload.Text,
	}
	modDataPayloadJson, err := json.Marshal(modDataPayload)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v: %w", err, asynq.SkipRetry)
	}

	modEvent := models.ModerationEvents{
		ContentId: payload.ContentID,
		EventType: models.EventType(models.Moderated),
		Payload:   modDataPayloadJson,
	}
	db.Create(&modEvent)

	status := models.ContentStatus(result.Status)
	task, err := tasks.NewAggregationDeliveryTask(payload.ContentID, &status, nil, nil)
	if err != nil {
		return fmt.Errorf("tasks.NewAggregationDeliveryTask failed: %v: %w", err, asynq.SkipRetry)
	}

	info, err := workerClient.Client.Enqueue(task, asynq.Queue(tasks.QueueAggregation))
	if err != nil {
		return fmt.Errorf("workerClient.Client.Enqueue failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	fmt.Println("Text processing completed")
	return nil
}
