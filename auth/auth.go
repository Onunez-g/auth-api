package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/gologin/github"
	"github.com/dghubble/sessions"
	"github.com/dgrijalva/jwt-go"
	"github.com/onunez-g/auth-api/config"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionName     = "auth-api"
	sessionUserKey  = "GithubID"
	sessionUsername = "githubUsername"
)

var sessionStore *sessions.CookieStore

func GenerateJWT(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	signingKey := config.Cfg.GetJWTSecret()

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		e := fmt.Errorf("Something went wrong: %s", err.Error())
		return "", e
	}

	return tokenString, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		signingKey := config.Cfg.GetJWTSecret()
		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return signingKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func IssueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		sessionStore = sessions.NewCookieStore(config.Cfg.GetSessionSecret(), nil)
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = *githubUser.ID
		session.Values[sessionUsername] = *githubUser.Login
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}
