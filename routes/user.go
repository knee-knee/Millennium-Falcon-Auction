package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/millennium-falcon-auction/repo"
)

// Login is the input object for when a user is logging in.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a user object.
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Session  string `json:"session,ommitempty"`
}

// Login will take in an email and password and return a valid session to the client.
func (r *Routes) Login(w http.ResponseWriter, req *http.Request) {
	var in Login
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		log.Printf("routes: Error trying to decode input %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	if in.Email == "" || in.Password == "" {
		log.Println("routes: email and password not provied for login request")
		http.Error(w, invalidUserNameOrPasswordResponse, http.StatusBadRequest)
		return
	}

	log.Printf("routes: Attempting to login user %s \n", in.Email)

	resp, err := r.Repo.GetUser(in.Email)
	if err != nil {
		log.Println("routes: Error trying to get user")
		http.Error(w, invalidUserNameOrPasswordResponse, http.StatusBadRequest)
		return
	}

	if resp.Password != in.Password {
		log.Printf("routes: User trying to login with invalid password %s \n", in.Email)
		http.Error(w, invalidUserNameOrPasswordResponse, http.StatusBadRequest)
		return
	}

	log.Printf("routes: Successfully logged in user %s \n", in.Email)

	w.Write([]byte(resp.Session))
}

// CreateUser will create a new user and return the session for that user to use.
func (r *Routes) CreateUser(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: Starting to create user")
	var in User
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		log.Printf("routes: Error trying to decode input %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()

	// Check if email is unique.
	userCheck, err := r.Repo.GetUser(in.Email)
	if err != nil {
		log.Printf("routes: could not get user %v", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	if (userCheck != repo.User{}) {
		log.Printf("routes: User tried to create a user with existing email %s \n", in.Email)
		http.Error(w, "user with email already exists", http.StatusBadRequest)
		return
	}

	// Create User
	user, err := r.Repo.CreateUser(in.Email, in.Password)
	if err != nil {
		log.Printf("routes: could not create user %v", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(user.Session))
}
