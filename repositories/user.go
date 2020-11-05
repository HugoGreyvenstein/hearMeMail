package repositories

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hearMeMail/global"
	"hearMeMail/models"
	"log"
	"time"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrDatabaseError = errors.New("database error")
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

func (database *UserRepository) FindByUsername(username string) (*models.User, error) {
	user := new(models.User)
	tx := database.connection.First(user, "username = ?", username)
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if tx.Error != nil {
		errMessage := ErrDatabaseError.Error() + fmt.Sprintf(": %+v", tx.Error)
		return user, errors.New(errMessage)
	}
	return user, nil
}

func (database *UserRepository) FindAll() []models.User {
	users := make([]models.User, 0)
	_ = database.connection.Find(&users, "username = ?", "goop")
	return users
}

func (database *UserRepository) Exists(user *models.User) (*models.User, error) {
	result := database.connection.First(user)
	return user, result.Error
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
	user, err := database.FindByUsername(user.Username)
	return user, err
}

func (database *UserRepository) DeleteHeaderToken(username string) error {
	expiry := time.Now()
	user := models.User{
		Username:     username,
		HeaderToken:  []byte{},
		HeaderExpiry: &expiry,
	}
	tx := database.connection.Model(&user).Updates(user)
	return tx.Error
}
