package repo

import "log"

type Repo struct{}

func New() *Repo {
	log.Println("instantiating a new repo")
	return &Repo{}
}
