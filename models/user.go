package models

import (
	"gorm.io/gorm"
	"time"
)

const EmailLogsAssociation = "EmailLogs"

var UserSchema = []interface{}{
	User{},
}

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Password     []byte
	HeaderToken  []byte
	HeaderExpiry *time.Time
	EmailLogs    []EmailLog
}
