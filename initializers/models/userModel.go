package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	First_name string
	Last_name  string
	Username   string
	Email      string `gorm:"unique"`
	Password   string
	Phone      string
	User_type  string `gorm:"default:'user'"`
	User_id    string `gorm:"unique"`
}
