package middleware

import "github.com/millennium-falcon-auction/repo"

// Middleware represents all shared values the middleware will need to be aware of.
type Middleware struct {
	Repo *repo.Repo
}

// New will return a new instance of the middleware object.
func New(r *repo.Repo) *Middleware {
	return &Middleware{
		Repo: r,
	}
}
