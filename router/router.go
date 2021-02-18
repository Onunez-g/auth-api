package router

import (
	"github.com/gorilla/mux"
	"github.com/onunez-g/auth-api/controllers"
)

func Get() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/auth/signup", controllers.SignUp).Methods("POST")
	r.HandleFunc("/api/auth/login", controllers.Login).Methods("POST")
	return r
}
