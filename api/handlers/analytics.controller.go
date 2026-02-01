package handlers

import (
	"net/http"
	"time"

	"github.com/Sreejit-Sengupto/internal/database"
	"github.com/Sreejit-Sengupto/internal/models"
	"github.com/Sreejit-Sengupto/utils/response"
)

type StatusCount struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type MediaTypeCount struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type RiskScoreRange struct {
	Range string `json:"range"`
	Count int64  `json:"count"`
}

type TimeSeriesData struct {
	Labels   []string         `json:"labels"`
	Datasets []DatasetEntry   `json:"datasets"`
}

type DatasetEntry struct {
	Label string  `json:"label"`
	Data  []int64 `json:"data"`
}

type RadarDataPoint struct {
	MediaType string `json:"mediaType"`
	Approved  int64  `json:"approved"`
	Rejected  int64  `json:"rejected"`
	Flagged   int64  `json:"flagged"`
	Pending   int64  `json:"pending"`
}

type ModerationSummary struct {
	TotalContent     int64 `json:"totalContent"`
	PendingCount     int64 `json:"pendingCount"`
	ApprovedCount    int64 `json:"approvedCount"`
	RejectedCount    int64 `json:"rejectedCount"`
	FlaggedCount     int64 `json:"flaggedCount"`
	TotalAudits      int64 `json:"totalAudits"`
	AvgRiskScore     float64 `json:"avgRiskScore"`
	ContentLastWeek  int64 `json:"contentLastWeek"`
}

func GetStatusDistribution(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	var results []StatusCount
	db.Model(&models.Content{}).
		Select("final_status as label, COUNT(*) as value").
		Group("final_status").
		Scan(&results)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": results,
	})
}

func GetModerationOverTime(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	type DailyCount struct {
		Date   string `json:"date"`
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var dailyCounts []DailyCount
	db.Model(&models.Content{}).
		Select("DATE(created_at) as date, final_status as status, COUNT(*) as count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Group("DATE(created_at), final_status").
		Order("date").
		Scan(&dailyCounts)

	dateMap := make(map[string]map[string]int64)
	var labels []string
	labelSet := make(map[string]bool)

	for _, dc := range dailyCounts {
		if !labelSet[dc.Date] {
			labels = append(labels, dc.Date)
			labelSet[dc.Date] = true
		}
		if dateMap[dc.Date] == nil {
			dateMap[dc.Date] = make(map[string]int64)
		}
		dateMap[dc.Date][dc.Status] = dc.Count
	}

	statuses := []string{"APPROVED", "REJECTED", "FLAGGED", "PENDING"}
	var datasets []DatasetEntry

	for _, status := range statuses {
		var data []int64
		for _, label := range labels {
			if val, ok := dateMap[label][status]; ok {
				data = append(data, val)
			} else {
				data = append(data, 0)
			}
		}
		datasets = append(datasets, DatasetEntry{
			Label: status,
			Data:  data,
		})
	}

	response.JSON(w, http.StatusOK, TimeSeriesData{
		Labels:   labels,
		Datasets: datasets,
	})
}

func GetMediaTypeBreakdown(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	var results []MediaTypeCount
	db.Model(&models.ModerationResult{}).
		Select("media_type as label, COUNT(*) as value").
		Group("media_type").
		Scan(&results)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": results,
	})
}

func GetRiskScoreDistribution(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	var results []RiskScoreRange

	ranges := []struct {
		Min   float64
		Max   float64
		Label string
	}{
		{0, 25, "0-25"},
		{25, 50, "26-50"},
		{50, 75, "51-75"},
		{75, 100, "76-100"},
	}

	for _, rng := range ranges {
		var count int64
		db.Model(&models.ModerationResult{}).
			Where("risk_score >= ? AND risk_score <= ?", rng.Min, rng.Max).
			Count(&count)

		results = append(results, RiskScoreRange{
			Range: rng.Label,
			Count: count,
		})
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": results,
	})
}

func GetStatusByMediaType(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	type MediaStatusCount struct {
		MediaType string `json:"mediaType"`
		Status    string `json:"status"`
		Count     int64  `json:"count"`
	}

	var counts []MediaStatusCount
	db.Model(&models.ModerationResult{}).
		Select("media_type, status, COUNT(*) as count").
		Group("media_type, status").
		Scan(&counts)

	mediaMap := make(map[string]*RadarDataPoint)
	for _, c := range counts {
		if mediaMap[c.MediaType] == nil {
			mediaMap[c.MediaType] = &RadarDataPoint{MediaType: c.MediaType}
		}
		switch c.Status {
		case "APPROVED":
			mediaMap[c.MediaType].Approved = c.Count
		case "REJECTED":
			mediaMap[c.MediaType].Rejected = c.Count
		case "FLAGGED":
			mediaMap[c.MediaType].Flagged = c.Count
		case "PENDING":
			mediaMap[c.MediaType].Pending = c.Count
		}
	}

	var results []RadarDataPoint
	for _, v := range mediaMap {
		results = append(results, *v)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": results,
	})
}

func GetAuditActivity(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	type DailyAudit struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	var dailyAudits []DailyAudit
	db.Model(&models.Audit{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -30)).
		Group("DATE(created_at)").
		Order("date").
		Scan(&dailyAudits)

	var labels []string
	var data []int64

	for _, da := range dailyAudits {
		labels = append(labels, da.Date)
		data = append(data, da.Count)
	}

	response.JSON(w, http.StatusOK, TimeSeriesData{
		Labels: labels,
		Datasets: []DatasetEntry{
			{Label: "Audits", Data: data},
		},
	})
}

func GetModerationSummary(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	var summary ModerationSummary

	db.Model(&models.Content{}).Count(&summary.TotalContent)
	db.Model(&models.Content{}).Where("final_status = ?", "PENDING").Count(&summary.PendingCount)
	db.Model(&models.Content{}).Where("final_status = ?", "APPROVED").Count(&summary.ApprovedCount)
	db.Model(&models.Content{}).Where("final_status = ?", "REJECTED").Count(&summary.RejectedCount)
	db.Model(&models.Content{}).Where("final_status = ?", "FLAGGED").Count(&summary.FlaggedCount)
	db.Model(&models.Audit{}).Count(&summary.TotalAudits)

	var avgScore struct {
		Avg float64
	}
	db.Model(&models.ModerationResult{}).Select("COALESCE(AVG(risk_score), 0) as avg").Scan(&avgScore)
	summary.AvgRiskScore = avgScore.Avg

	db.Model(&models.Content{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Count(&summary.ContentLastWeek)

	response.JSON(w, http.StatusOK, summary)
}
