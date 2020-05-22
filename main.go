package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/millennium-falcon-auction/middleware"
	"github.com/millennium-falcon-auction/repo"
	"github.com/millennium-falcon-auction/routes"
)

func main() {
	log.Println("Now strarting the web server")

	repo := repo.New()
	routes := routes.New(repo)
	mw := middleware.New(repo)

	r := mux.NewRouter()

	// health check
	r.HandleFunc("/", routes.HealthCheck).Methods(http.MethodGet)

	// item info
	r.HandleFunc("/item/{id}", routes.GetItemInfo).Methods(http.MethodGet)

	// bid
	r.Handle("/bid/{item_id}", mw.AuthMiddleware(routes.PlaceBid)).Methods(http.MethodPost)
	r.Handle("/bid/{id}", mw.AuthMiddleware(routes.UpdateBid)).Methods(http.MethodPatch)

	// User
	r.HandleFunc("/user/login", routes.Login).Methods(http.MethodPost)
	r.HandleFunc("/user/signup", routes.CreateUser).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}
