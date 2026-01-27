package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res, _ := json.Marshal("All ok")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
