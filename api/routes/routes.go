package routes

import "github.com/gorilla/mux"

func RegisterRoutes(r *mux.Router) {
	registerUploadRoutes(r)
	registerContentRoutes(r)
	registerTestRoutes(r)
}
