package routes

import (
	"github.com/Sreejit-Sengupto/api/handlers"
	"github.com/gorilla/mux"
)

func registerUploadRoutes(r *mux.Router) {
	r.HandleFunc("/upload/content", handlers.UploadContent).Methods("POST")
}
