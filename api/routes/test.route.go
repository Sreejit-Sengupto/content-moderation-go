package routes

import (
	"github.com/Sreejit-Sengupto/api/handlers"
	"github.com/gorilla/mux"
)

func registerTestRoutes(r *mux.Router) {
	r.HandleFunc("/test", handlers.HealthTest).Methods("GET")
}
