package routes

import (
	"fmt"
	"log"
	"net/http"
)

func (r *Routes) HealthCheck(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Running a health check")
	fmt.Fprintf(w, "healthy")
}
