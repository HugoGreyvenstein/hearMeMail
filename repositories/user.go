package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hearMeMail/global"
	"hearMeMail/models"
	"log"
	"time"
)

type UserRepository struct {
	config     *global.Config
	connection *gorm.DB
}

func UserRepositoryBuild(config *global.Config) *UserRepository {
	return &UserRepository{config: config}
}

func (database *UserRepository) Initialise() error {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable",
		database.config.Database.Credentials.Username,
		database.config.Database.Credentials.Password,
		database.config.Database.Name,
		database.config.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error initialising database: err=%+v", err)
		return err
	}
	database.connection = db

	err = database.connection.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("Failed to migrate schema 'User': err=%+v", err)
	}
	return nil
}

func (database *UserRepository) FindByUsername(username string) *models.User {
	user := new(models.User)
	database.connection.First(user, "username = ?", username)
	if user.Username == "" {
		return nil
	}
	return user
}

func (database *UserRepository) FindAll() []models.User {
	users := make([]models.User, 0)
	_ = database.connection.Find(&users, "username = ?", "goop")
	return users
}

func (database *UserRepository) Insert(user *models.User) (*models.User, error) {
	result := database.connection.FirstOrCreate(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (database *UserRepository) UpdateHeaderToken(user *models.User, headerToken []byte, expiry time.Time) (*models.User, error) {
	result := database.connection.Model(&user).Updates(models.User{
		HeaderToken:  headerToken,
		HeaderExpiry: &expiry,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	return database.FindByUsername(user.Username), nil
}

func (database *UserRepository) DeleteHeaderToken(username string, token string) error {
	user := models.User{
		Username:    username,
		HeaderToken: []byte(token),
	}
	tx := database.connection.Model(&user).Updates(models.User{HeaderToken: nil, HeaderExpiry: nil})
	return tx.Error
}
