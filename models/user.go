package models

import "gorm.io/gorm"

type User struct {
	model    gorm.Model
	ID       int    `json:"id" gorm:"column:user_id;primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"column:user_name"`
	Email    string `json:"email" gorm:"column:email"`
	Phone    string `json:"phone_number" gorm:"column:phone_number"`
	Password string `json:"password" gorm:"column:password"`
	Role_id  int    `json:"role_id" gorm:"column:role_id"`
}

type APIUser struct {
	Id   int    `json:"id" gorm:"column:user_id" `
	Name string `json:"name" gorm:"column:user_name"`
}
