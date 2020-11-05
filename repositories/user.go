package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hearMeMail/global"
	"hearMeMail/models"
	"time"
)

type UserRepository struct {
	config     *global.Config
	connection *gorm.DB
}

func UserRepositoryBuild(config *global.Config) *UserRepository {
	return &UserRepository{config: config}
}

func (database *UserRepository) InitialiseConnection() error {
	conn, err := getConnection(database.config)
	database.connection = conn
	return err
}

func (database *UserRepository) AutoMigrate() {
	autoMigrate(database.connection, models.UserSchema)
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
	_ = database.connection.Find(&users)
	return users
}

func (database *UserRepository) Exists(user *models.User) (*models.User, error) {
	result := database.connection.First(user)
	return user, result.Error
}

func (database *UserRepository) Insert(user *models.User) (*models.User, error) {
	existingUser := new(models.User)
	tx := database.connection.Find(existingUser)
	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
		return nil, tx.Error
	}
	result := database.connection.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return database.FindByUsername(user.Username)
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
