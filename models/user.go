package models

import "time"

type User struct {
	Username     string `gorm:"primaryKey"`
	Password     []byte
	HeaderToken  []byte
	HeaderExpiry *time.Time
}
