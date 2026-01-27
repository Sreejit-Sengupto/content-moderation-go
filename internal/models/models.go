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
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Text             string
	Image            string
	Video            string
	TextStatus       ContentStatus
	ImageStatus      ContentStatus
	VideoStatus      ContentStatus
	FinalStatus      ContentStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ModerationResult []ModerationResult
	ModerationEvents []ModerationEvents
	Audit            []Audit
}

type ModerationResult struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ContentId    uuid.UUID     `gorm:"not null"`
	Content      Content       `gorm:"foreignKey:ContentId"`
	MediaType    MediaType     `gorm:"not null"`
	Status       ContentStatus `gorm:"not null"`
	RiskScore    float64       `gorm:"not null"`
	Explaination string
	CreatedAt    time.Time
}

type ModerationEvents struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ContentId uuid.UUID      `gorm:"not null"`
	Content   Content        `gorm:"foreignKey:ContentId"`
	EventType EventType      `gorm:"not null"`
	Payload   datatypes.JSON `gorm:"type:JSONB"`
	CreatedAt time.Time
}

type Audit struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ContentId uuid.UUID `gorm:"not null"`
	Content   Content   `gorm:"foreignKey:ContentId"`
	Action    Action    `gorm:"not null"`
	Reason    string
	CreatedAt time.Time
}
