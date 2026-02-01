package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sreejit-Sengupto/internal/database"
	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/Sreejit-Sengupto/utils/response"
	"github.com/Sreejit-Sengupto/utils/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAllContent(w http.ResponseWriter, r *http.Request) {
	var contents []models.Content
	db := database.DB
	result := db.Find(&contents)
	if result.Error != nil {
		response.JSONError(w, http.StatusNotFound, "Failed to fetch all content")
	}
	response.JSON(w, http.StatusOK, contents)
}

func GetContentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "Invalid content ID")
		return
	}

	db := database.DB
	var content models.Content
	result := db.Find(&content, models.Content{ID: id})
	if result.Error != nil {
		response.JSONError(w, http.StatusNotFound, "Failed to fetch content details for the id "+idStr)
	}
	response.JSON(w, http.StatusOK, content)
}

func GetModerationResults(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "Invalid content ID")
		return
	}

	db := database.DB

	var content models.Content

	result := db.Preload("ModerationResult").Find(&content, models.Content{ID: id})
	if result.Error != nil {
		fmt.Println(result.Error)
		response.JSONError(w, http.StatusNotFound, "Failed to fetch results")
		return
	}

	response.JSON(w, http.StatusOK, content)
}

func GetModerationEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "Invalid content ID")
		return
	}

	db := database.DB

	var content models.Content
	result := db.Preload("ModerationEvents").Find(&content, models.Content{ID: id})
	if result.Error != nil {
		response.JSONError(w, http.StatusNotFound, "Failed to fetch events")
		return
	}
	response.JSON(w, http.StatusOK, content)
}

func GetModerationAudits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.JSONError(w, http.StatusBadRequest, "Invalid content ID")
		return
	}

	db := database.DB

	var content models.Content
	result := db.Preload("Audit").Find(&content, models.Content{ID: id})
	if result.Error != nil {
		response.JSONError(w, http.StatusNotFound, "Failed to fetch events")
		return
	}
	response.JSON(w, http.StatusOK, content)
}

func UpdateContent(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		ContentID   uuid.UUID `json:"contentId" validate:"required"`
		Reason      string    `json:"reason" validate:"required"`
		TextStatus  string    `json:"textStatus" validate:"oneof=PENDING APPROVED REJECTED FLAGGED"`
		ImageStatus string    `json:"imageStatus" validate:"oneof=PENDING APPROVED REJECTED FLAGGED"`
		VideoStatus string    `json:"videoStatus" validate:"oneof=PENDING APPROVED REJECTED FLAGGED"`
		FinalStatus string    `json:"finalStatus" validate:"oneof=PENDING APPROVED REJECTED FLAGGED"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		response.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validator.Validtor().Struct(reqBody); err != nil {
		response.JSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %v", err))
		return
	}

	db := database.DB

	var content models.Content

	result := db.Find(&content, models.Content{ID: reqBody.ContentID})
	if result.Error != nil {
		response.JSONError(w, http.StatusNotFound, "Failed to fetch all content")
	}
	content.TextStatus = models.ContentStatus(reqBody.TextStatus)
	content.ImageStatus = models.ContentStatus(reqBody.ImageStatus)
	content.VideoStatus = models.ContentStatus(reqBody.VideoStatus)
	content.FinalStatus = models.ContentStatus(reqBody.FinalStatus)

	db.Save(&content)

	auditLogs := models.Audit{
		ContentId: reqBody.ContentID,
		Action:    "OVERIDDEN",
		Reason:    reqBody.Reason,
	}
	result = db.Create(&auditLogs)

	if result.Error != nil {
		response.JSONError(w, http.StatusBadRequest, "Failed to insert audit logs")
		return
	}

	response.JSON(w, http.StatusCreated, "Content updated and audited")
}
