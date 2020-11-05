package models

import "gorm.io/gorm"

var EmailLogSchema = []interface{}{
	EmailLog{},
}

type EmailLog struct {
	gorm.Model
	Subject string
	Body    string
	To      string
	UserID  uint
	Success bool
}
