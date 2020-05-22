package middleware

import "github.com/millennium-falcon-auction/repo"

type Middleware struct {
	Repo *repo.Repo
}

func New(r *repo.Repo) *Middleware {
	return &Middleware{
		Repo: r,
	}
}
