package data

import (
	"github.com/onunez-g/auth-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//Db database
var Db *gorm.DB
var err error

func ConnectDatabase() {
	Db, err = gorm.Open(sqlite.Open("auth-api.db"), &gorm.Config{})

	if err != nil {

		panic(err.Error())

	}
}

func AutoMigrate() {
	Db.AutoMigrate(&models.UserDTO{})
}

// func CloseConnection() {
// 	Db.Debug().Statement.ReflectValue.Close()
// }
