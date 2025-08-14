package models

import "gorm.io/gorm"

type Device struct {
	model       gorm.Model
	ID          int    `json:"device_id" gorm:"column:device_id"` // Device ID field
	Model       string `json:"model" gorm:"column:model"`         // Device model
	Brand       string `json:"brand" gorm:"column:brand"`         // Device brand
	UserId      int    `json:"user_id" gorm:"column:user_id"`     // User ID field
	Description string `json:"description" gorm:"column:description"`
}

type ApiDevice struct {
	model gorm.Model
	ID    int    `json:"device_id" gorm:"column:device_id"`
	Model string `json:"model" gorm:"column:model"`
}
