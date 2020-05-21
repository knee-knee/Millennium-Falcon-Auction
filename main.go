package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/millennium-falcon-auction/repo"
	"github.com/millennium-falcon-auction/routes"
)

func main() {
	log.Println("Now strarting the web server")

	repo := repo.New()
	routes := routes.New(repo)

	r := mux.NewRouter()

	// health check
	r.HandleFunc("/", routes.HealthCheck).Methods(http.MethodGet)

	// item info
	r.HandleFunc("/item/{id}", routes.GetItemInfo).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
