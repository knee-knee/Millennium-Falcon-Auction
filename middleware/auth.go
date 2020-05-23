package middleware

import (
	"log"
	"net/http"

	"github.com/millennium-falcon-auction/repo"
)

// AuthHeader is the key for the auth header.
const AuthHeader = "auth"

// AuthMiddleware is the handler for the auth middleware.AuthMiddleware
// AuthMiddleware will ensure that the a proper session is provided for the auth header.
func (m *Middleware) AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("middleware: Trying to find user")
		token := req.Header.Get(AuthHeader)
		if token == "" {
			log.Println("mw: auth token was not provided for request")
			http.Error(w, "Did not provide auth header", http.StatusBadRequest)
			return
		}

		user, err := m.Repo.GetUserBySession(token)
		if err != nil {
			log.Println("mw: error trying to get the user by their session")
			http.Error(w, "Could not retrieve user from dynamo", http.StatusInternalServerError)
			return
		}

		if (user != repo.User{}) {
			log.Printf("middleware: User %s successfully authed. \n", token)
			next.ServeHTTP(w, req)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
