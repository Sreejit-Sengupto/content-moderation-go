package routes

import (
	"github.com/Sreejit-Sengupto/api/handlers"
	"github.com/gorilla/mux"
)

func registerContentRoutes(r *mux.Router) {
	r.HandleFunc("/content", handlers.GetAllContent).Methods("GET", "OPTIONS")
	r.HandleFunc("/content/update", handlers.UpdateContent).Methods("PATCH", "OPTIONS")
	r.HandleFunc("/content/{id}", handlers.GetContentByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/content/{id}/results", handlers.GetModerationResults).Methods("GET", "OPTIONS")
	r.HandleFunc("/content/{id}/events", handlers.GetModerationEvents).Methods("GET", "OPTIONS")
	r.HandleFunc("/content/{id}/audits", handlers.GetModerationAudits).Methods("GET", "OPTIONS")
}
