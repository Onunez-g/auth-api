package router

import (
	"github.com/gorilla/mux"
	"github.com/onunez-g/auth-api/controllers"
)

func Get() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/api/auth/signup", controllers.SignUp).Methods("POST")
	return r
}
