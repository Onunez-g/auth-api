package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/onunez-g/auth-api/config"
	"github.com/onunez-g/auth-api/data"
	"github.com/onunez-g/auth-api/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Failed to load env vars")
	}
	cfg := config.Get()
	data.ConnectDatabase()

	r := router.Get()

	log.Println("Server listening...")
	apiPort := cfg.GetAPIPort()

	log.Fatal(http.ListenAndServe(apiPort, r))
}
