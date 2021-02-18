package data

import (
	"log"

	"github.com/onunez-g/auth-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Db database
var Db *gorm.DB
var err error

func ConnectDatabase(connStr string) {
	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {

		panic(err.Error())

	}
}

func AutoMigrate() {
	err = Db.AutoMigrate(&models.UserDTO{})
	if err != nil {
		log.Println("something happenned in the migration:" + err.Error())
	}
}

func CloseConnection() {
	db, err := Db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	db.Close()
}
