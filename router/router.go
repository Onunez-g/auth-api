package router

import (
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/gorilla/mux"
	"github.com/onunez-g/auth-api/auth"
	"github.com/onunez-g/auth-api/config"
	"github.com/onunez-g/auth-api/controllers"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

func Get() *mux.Router {
	r := mux.NewRouter()

	stateConfig := gologin.DebugOnlyCookieConfig

	oauth2Config := &oauth2.Config{
		ClientID:     config.Cfg.GetGithubClientId(),
		ClientSecret: config.Cfg.GetGithubSecretId(),
		RedirectURL:  "http://localhost:8080/api/github/callback",
		Endpoint:     githubOAuth2.Endpoint,
	}

	r.HandleFunc("/api/auth/signup", controllers.SignUp).Methods("POST")
	r.HandleFunc("/api/auth/login", controllers.Login).Methods("POST")
	r.HandleFunc("/api/auth/sendemail/{email}", controllers.SendConfirmationEmail).Methods("GET")
	r.Handle("/api/github/login", github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil)))
	r.Handle("/api/github/callback", github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, auth.IssueSession(), nil)))
	return r
}
