package routes

import (
	"github.com/Sreejit-Sengupto/api/handlers"
	"github.com/gorilla/mux"
)

func registerAnalyticsRoutes(r *mux.Router) {
	r.HandleFunc("/analytics/status-distribution", handlers.GetStatusDistribution).Methods("GET", "OPTIONS")
	r.HandleFunc("/analytics/moderation-over-time", handlers.GetModerationOverTime).Methods("GET", "OPTIONS")
	r.HandleFunc("/analytics/media-type-breakdown", handlers.GetMediaTypeBreakdown).Methods("GET", "OPTIONS")
	r.HandleFunc("/analytics/risk-score-distribution", handlers.GetRiskScoreDistribution).Methods("GET", "OPTIONS")
	r.HandleFunc("/analytics/status-by-media-type", handlers.GetStatusByMediaType).Methods("GET", "OPTIONS")
	r.HandleFunc("/analytics/audit-activity", handlers.GetAuditActivity).Methods("GET", "OPTIONS")
	r.HandleFunc("/analytics/summary", handlers.GetModerationSummary).Methods("GET", "OPTIONS")
}
