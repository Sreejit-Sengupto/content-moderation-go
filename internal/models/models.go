package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// 'PENDING', 'APPROVED', 'REJECTED', 'FLAGGED'
type ContentStatus string

// 'TXT', 'IMG', 'VIDEO'
type MediaType string

// 'CREATED', 'UPDATED', 'MODERATED'
type EventType string

// 'REVIEWED', 'OVERRIDEN'
type Action string

const (
	Pending  ContentStatus = "PENDING"
	Approved ContentStatus = "APPROVED"
	Rejected ContentStatus = "REJECTED"
	Flagged  ContentStatus = "FLAGGED"
)

const (
	Txt MediaType = "TXT"
	Img MediaType = "IMG"
	Vid MediaType = "VID"
)

const (
	Created   EventType = "CREATED"
	Updated   EventType = "UPDATED"
	Moderated EventType = "MODERATED"
)

const (
	Reviewed  Action = "REVIEWED"
	Overriden Action = "OVERRIDEN"
)

type Content struct {
	ID               uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Text             string             `json:"text"`
	Image            string             `json:"image"`
	Video            string             `json:"video"`
	TextStatus       ContentStatus      `json:"textStatus"`
	ImageStatus      ContentStatus      `json:"imageStatus"`
	VideoStatus      ContentStatus      `json:"videoStatus"`
	FinalStatus      ContentStatus      `json:"finalStatus"`
	CreatedAt        time.Time          `json:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt"`
	ModerationResult []ModerationResult `json:"moderationResult,omitempty"`
	ModerationEvents []ModerationEvents `json:"moderationEvents,omitempty"`
	Audit            []Audit            `json:"audits"`
}

type ModerationResult struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ContentId    uuid.UUID     `gorm:"not null" json:"contentId"`
	Content      Content       `gorm:"foreignKey:ContentId" json:"-"`
	MediaType    MediaType     `gorm:"not null" json:"mediaType"`
	Status       ContentStatus `gorm:"not null" json:"status"`
	RiskScore    float64       `gorm:"not null" json:"riskScore"`
	Explaination string        `json:"explanation"`
	CreatedAt    time.Time     `json:"createdAt"`
}

type ModerationEvents struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ContentId uuid.UUID      `gorm:"not null" json:"contentId"`
	Content   Content        `gorm:"foreignKey:ContentId" json:"-"`
	EventType EventType      `gorm:"not null" json:"eventType"`
	Payload   datatypes.JSON `gorm:"type:JSONB" json:"payload"`
	CreatedAt time.Time      `json:"createdAt"`
}

type Audit struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ContentId uuid.UUID `gorm:"not null" json:"contentId"`
	Content   Content   `gorm:"foreignKey:ContentId" json:"-"`
	Action    Action    `gorm:"not null" json:"action"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
}
