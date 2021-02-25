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

	apiPort := config.Cfg.GetAPIPort()
	log.Println("Server listening at port: " + apiPort)

	log.Fatal(http.ListenAndServe(apiPort, r))
}
