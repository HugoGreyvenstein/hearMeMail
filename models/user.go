package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username     string `gorm:"primaryKey"`
	Password     []byte
	HeaderToken  []byte
	HeaderExpiry *time.Time
}
