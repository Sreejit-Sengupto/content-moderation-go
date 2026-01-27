package aggregation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sreejit-Sengupto/internal/database"
	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/Sreejit-Sengupto/internal/queue/tasks"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func HandleAggregationDelivery(ctx context.Context, t *asynq.Task) error {
	fmt.Println("Processing aggregation delivery task")
	var payload tasks.ResultAggregationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	db := database.DB

	// Use transaction with row-level locking to prevent race conditions
	// when multiple aggregation tasks run concurrently for the same content
	err := db.Transaction(func(tx *gorm.DB) error {
		var existingContent models.Content

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(&models.Content{ID: payload.ContentID}).
			First(&existingContent).Error; err != nil {
			return fmt.Errorf("failed to find content: %v", err)
		}

		if existingContent.TextStatus == "" && payload.TextStatus != nil {
			existingContent.TextStatus = *payload.TextStatus
		}
		if existingContent.ImageStatus == "" && payload.ImageStatus != nil {
			existingContent.ImageStatus = *payload.ImageStatus
		}
		if existingContent.VideoStatus == "" && payload.VideoStatus != nil {
			existingContent.VideoStatus = *payload.VideoStatus
		}

		// Priority:
		// REJECTED > FLAGGED > PENDING > APPROVED
		// - if any one is rejected -> final is rejected
		// - otherwise if two of them are flagged -> final is flagged
		// - approve only if all non-empty are approved
		// - otherwise pending

		var allStatuses []models.ContentStatus

		if existingContent.TextStatus != "" {
			allStatuses = append(allStatuses, existingContent.TextStatus)
		}
		if existingContent.ImageStatus != "" {
			allStatuses = append(allStatuses, existingContent.ImageStatus)
		}
		if existingContent.VideoStatus != "" {
			allStatuses = append(allStatuses, existingContent.VideoStatus)
		}

		var finalStatus models.ContentStatus

		if len(allStatuses) == 0 {
			finalStatus = models.Pending
		} else {
			hasRejected := false
			flaggedCount := 0
			allApproved := true

			for _, s := range allStatuses {
				if s == models.Rejected {
					hasRejected = true
					break
				}
				if s == models.Flagged {
					flaggedCount++
				}
				if s != models.Approved {
					allApproved = false
				}
			}

			if hasRejected {
				finalStatus = models.Rejected
			} else if flaggedCount >= 2 {
				finalStatus = models.Flagged
			} else if allApproved {
				finalStatus = models.Approved
			} else {
				finalStatus = models.Pending
			}
		}

		existingContent.FinalStatus = finalStatus
		if err := tx.Save(&existingContent).Error; err != nil {
			return fmt.Errorf("failed to update content: %v", err)
		}

		fmt.Printf("Aggregation complete for content %s: final status = %s\n", payload.ContentID, finalStatus)
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %v", err)
	}

	return nil
}
