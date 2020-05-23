package routes

import (
	"log"

	"github.com/millennium-falcon-auction/repo"
)

const (
	invalidUserNameOrPasswordResponse = "Invalid Email or Password"
	internalErrorResponse             = "Internal Server Error"
)

type Routes struct {
	Repo *repo.Repo
}

func New(r *repo.Repo) *Routes {
	log.Println("Instantiating a new route")
	return &Routes{
		Repo: r,
	}
}
