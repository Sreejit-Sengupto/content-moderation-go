package routes

import (
	"github.com/Sreejit-Sengupto/api/handlers"
	"github.com/gorilla/mux"
)

func registerContentRoutes(r *mux.Router) {
	r.HandleFunc("/content", handlers.GetAllContent).Methods("GET", "OPTIONS")
	r.HandleFunc("/content/update", handlers.UpdateContent).Methods("PATCH", "OPTIONS")
}
