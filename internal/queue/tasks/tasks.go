package tasks

import (
	"encoding/json"

	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// Task types (used for routing to handlers)
const (
	TypeTextDelivery        = "text_delivery"
	TypeImageDelivery       = "image_delivery"
	TypeVideoDelivery       = "video_delivery"
	TypeAggregationDelivery = "aggregation"
)

// Queue names (used for queue assignment)
const (
	QueueText        = "text"
	QueueImage       = "image"
	QueueVideo       = "video"
	QueueAggregation = "aggregation"
)

type TextDeliveryPayload struct {
	ContentID uuid.UUID
	Text      string
}

type ImageDeliveryPayload struct {
	ContentID uuid.UUID
	Image     string
}

type VideoDeliveryPayload struct {
	ContentID uuid.UUID
	Video     string
}

type ResultAggregationPayload struct {
	ContentID   uuid.UUID
	TextStatus  *models.ContentStatus
	ImageStatus *models.ContentStatus
	VideoStatus *models.ContentStatus
}

func NewTextDeliveryTask(contentId uuid.UUID, text string) (*asynq.Task, error) {
	payload, err := json.Marshal(TextDeliveryPayload{
		ContentID: contentId,
		Text:      text,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeTextDelivery, payload), nil
}

func NewImageDeliveryTask(contentId uuid.UUID, image string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageDeliveryPayload{
		ContentID: contentId,
		Image:     image,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeImageDelivery, payload), nil
}

func NewVideoDeliveryTask(contentId uuid.UUID, video string) (*asynq.Task, error) {
	payload, err := json.Marshal(VideoDeliveryPayload{
		ContentID: contentId,
		Video:     video,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeVideoDelivery, payload), nil
}

func NewAggregationDeliveryTask(contentId uuid.UUID, textStatus *models.ContentStatus, imageStatus *models.ContentStatus, videoStatus *models.ContentStatus) (*asynq.Task, error) {
	payload, err := json.Marshal(ResultAggregationPayload{
		ContentID:   contentId,
		TextStatus:  textStatus,
		ImageStatus: imageStatus,
		VideoStatus: videoStatus,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeAggregationDelivery, payload), nil
}
