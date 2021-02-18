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
	config.Cfg = config.Get()

	data.ConnectDatabase(config.Cfg.GetDBConnStr())
	defer data.CloseConnection()

	r := router.Get()

	data.AutoMigrate()

	log.Println("Server listening...")
	apiPort := config.Cfg.GetAPIPort()

	log.Fatal(http.ListenAndServe(apiPort, r))
}
