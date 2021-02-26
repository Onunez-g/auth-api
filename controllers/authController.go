package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onunez-g/auth-api/auth"
	"github.com/onunez-g/auth-api/config"
	"github.com/onunez-g/auth-api/data"
	"github.com/onunez-g/auth-api/mail"
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

	err = sendEmail(user.Email, "Welcome to GoAuth!", "Hello there new fella!")
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
	err := sendEmail(params["email"], "Welcome to GoAuth!", "Hello there new fella!")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to send email, retry later"))
		return
	}
	w.Write([]byte("Email sent"))
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

	jwtToken, err := auth.GenerateJWT(user.Username)
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

func sendEmail(to string, subject string, body string) error {
	email := mail.EmailServer{
		User: config.Cfg.GetSMTPUser(),
		Pass: config.Cfg.GetSMTPPassword(),
		From: config.Cfg.GetSMTPUser(),
		Smtp: "smtp.gmail.com",
		Port: 587,
		To:   to,
	}
	return email.Send(subject, body)
}
