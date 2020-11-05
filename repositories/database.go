package repositories

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hearMeMail/global"
	"hearMeMail/models"
	"log"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrDatabaseError = errors.New("database error")
)

func autoMigrate(connection *gorm.DB, schemas []interface{}) {
	for _, schema := range schemas {
		err := connection.AutoMigrate(schema)
		if err != nil {
			log.Printf("Migration failed: %+v", err)
		}
	}
}

func getConnection(config *global.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Database.Credentials.Username,
		config.Database.Credentials.Password,
		config.Database.Name,
		config.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error initialising database: err=%+v", err)
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("Failed to migrate schema 'User': err=%+v", err)
	}
	return db, nil
}
