package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/onunez-g/auth-api/auth"
	"github.com/onunez-g/auth-api/data"
	"github.com/onunez-g/auth-api/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	var user models.UserDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte(err.Error()))
	}
	user.Password = auth.GetHash([]byte(user.Password))

	err = data.Db.Create(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

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
