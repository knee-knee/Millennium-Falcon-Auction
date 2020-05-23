package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

// HealthCheckOutput is the output when a health check is preformed.
type HealthCheckOutput struct {
	Status string `json:"status"`
}

// HealthCheck will peform a healthcheck to make sure the server is okay.
func (r *Routes) HealthCheck(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Running a health check")
	out := HealthCheckOutput{
		Status: "Healthy",
	}
	body, err := json.Marshal(out)
	if err != nil {
		log.Printf("routes: Error marshalling in to response body %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
