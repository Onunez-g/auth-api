package models

type UserDTO struct {
	Id       int    `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
}
