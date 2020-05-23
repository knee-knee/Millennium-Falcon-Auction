package routes

import (
	"log"
	"sync"

	"github.com/millennium-falcon-auction/repo"
)

const (
	invalidUserNameOrPasswordResponse = "Invalid Email or Password"
	internalErrorResponse             = "Internal Server Error"
)

// Routes is the shared information that is need for the routes functions.
type Routes struct {
	Repo          *repo.Repo
	highestBidMux sync.Mutex
}

// New will return a new instance of the routes object.
func New(r *repo.Repo) *Routes {
	log.Println("Instantiating a new route")
	return &Routes{
		Repo: r,
	}
}
