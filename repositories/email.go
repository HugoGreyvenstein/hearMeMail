package repositories

import (
	"gorm.io/gorm"
	"hearMeMail/global"
	"hearMeMail/models"
)

type EmailLogRepository struct {
	config     *global.Config
	connection *gorm.DB
}

func EmailLogRepositoryBuild(config *global.Config) *EmailLogRepository {
	return &EmailLogRepository{config: config}
}

func (database *EmailLogRepository) InitialiseConnection() error {
	conn, err := getConnection(database.config)
	database.connection = conn
	return err
}

func (database *EmailLogRepository) AutoMigrate() {
	autoMigrate(database.connection, models.EmailLogSchema)
}

func (database *EmailLogRepository) FindAllByUser(user *models.User) ([]models.EmailLog, error) {
	emailLogs := make([]models.EmailLog, 0)
	err := database.connection.Model(&user).
		Association(models.EmailLogsAssociation).
		Find(emailLogs)
	return emailLogs, err
}

func (database *EmailLogRepository) LogEmail(user *models.User, email *models.EmailLog) error {
	return database.connection.Model(&user).
		Association(models.EmailLogsAssociation).
		Append(email)
}
