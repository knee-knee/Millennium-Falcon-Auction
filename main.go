package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/millennium-falcon-auction/routes"
)

func main() {
	log.Println("Now strarting the web server")
	r := mux.NewRouter()
	// health check
	r.HandleFunc("/", routes.HealthCheck).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
