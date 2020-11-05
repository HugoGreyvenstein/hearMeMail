package main

import (
	"hearMeMail/global"
	"hearMeMail/handlers"
	"hearMeMail/repositories"
	"hearMeMail/services"
	"log"
	"net/http"
	"os"
)

const configFile = "config.yml"

func main() {
	configFileName := configFile
	args := os.Args
	if len(args) > 1 {
		configFileName = args[1]
	}
	config, err := global.LoadConfig(configFileName)
	if err != nil {
		log.Printf("Error occurred loading config: %+v", err)
		err = nil
	}

	// Create Repositories
	userRepository := repositories.UserRepositoryBuild(config)
	err = userRepository.InitialiseConnection()
	if err != nil {
		log.Printf("User repository initialisation failure: %+v", err)
		err = nil
	}
	userRepository.AutoMigrate()

	emailLogRepository := repositories.EmailLogRepositoryBuild(config)
	err = emailLogRepository.InitialiseConnection()
	if err != nil {
		log.Printf("EmailLogs repository initialisation failure: %+v", err)
		err = nil
	}
	emailLogRepository.AutoMigrate()

	// Create Services
	emailService := services.EmailServiceBuild(config, emailLogRepository, userRepository)
	loginService := services.LoginServiceBuild(config, userRepository)

	// Create Handlers
	emailHandler := handlers.EmailHandlerBuild(config, emailService)
	loginHandler := handlers.LoginHandlerBuild(config, loginService)
	logoutHandler := handlers.LogoutHandlerBuild(config, loginService)
	registerHandler := handlers.RegisterHandlerBuild(config, userRepository)

	// Register Handlers
	http.HandleFunc("/email", loginHandler.TokenChecker(emailHandler.Handler))
	http.HandleFunc("/login", loginHandler.Handler)
	http.HandleFunc("/logout", loginHandler.TokenChecker(logoutHandler.Handler))
	http.HandleFunc("/register", registerHandler.Handler)

	log.Print("Starting email server")
	// TODO Make port configurable
	err = http.ListenAndServe(":8080", nil)
	log.Printf("Error occurred while running server: err=%+v", err)
}
