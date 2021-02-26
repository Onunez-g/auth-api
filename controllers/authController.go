package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onunez-g/auth-api/auth"
	"github.com/onunez-g/auth-api/config"
	"github.com/onunez-g/auth-api/data"
	"github.com/onunez-g/auth-api/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	var user models.UserDTO
	var duplicatedUser models.UserDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte(err.Error()))
	}
	user.Password = auth.GetHash([]byte(user.Password))

	data.Db.Find(&duplicatedUser, "Email = ?", user.Email)
	if duplicatedUser.Email == user.Email {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("User already created"))
		return
	}
	err = data.Db.Create(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	mailSettings := config.Cfg.GetSMTPSettings(user.Email)
	err = mailSettings.Send("Welcome to GoAuth", "Hello there new fella!")
	if err != nil {
		log.Println(err.Error())
		log.Println("Unable to send email, retry later")
	}

	response := getResponse(&user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func SendConfirmationEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	payload := map[string]interface{}{
		"confirmEmail": false,
	}

	token, err := auth.GenerateJWT(params["email"], payload)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to generate token"))
	}
	mailSettings := config.Cfg.GetSMTPSettings(params["email"])
	body := "Hello there new fella! " + "Please confirm your email with this token: \n" + token
	err = mailSettings.Send("Welcome to GoAuth", body)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to send email, retry later"))
		return
	}
	w.Write([]byte("Email sent succesfully\n token: " + token))
}

func ActivateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.UserDTO
	claims, err := auth.GetJWT(params["token"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Errorf("Something went wrong: %s", err.Error())
		w.Write([]byte("There was an unexpected error"))
	}
	if claims["user"] == "" {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Errorf("Something went wrong: %s", err.Error())
		w.Write([]byte("There's no user assigned"))
	}
	data.Db.Find(&user, "Email = ?", claims["user"])
	user.EmailConfirmed = true
	data.Db.Updates(&user)

	response := getResponse(&user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.UserDTO
	var loginUser models.UserDTO
	json.NewDecoder(r.Body).Decode(&user)

	data.Db.Find(&loginUser, "Email = ?", user.Email)

	if loginUser.Id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect credentials"))
		return
	}
	userPass := []byte(user.Password)
	loginPass := []byte(loginUser.Password)

	passErr := bcrypt.CompareHashAndPassword(loginPass, userPass)

	if passErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect credentials"))
		return
	}
	payload := map[string]interface{}{
		"emailConfirmed": loginUser.EmailConfirmed,
	}
	jwtToken, err := auth.GenerateJWT(user.Username, payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	loginUser.Token = jwtToken

	response := getResponse(&loginUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getResponse(o interface{}) []byte {
	response, err := json.Marshal(&o)

	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return response
}
