package main

import (
	"hearMeMail/global"
	"hearMeMail/handlers"
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
	log.Printf("config=%+v", config)

	emailService := services.EmailServiceBuild(config)

	emailHandler := handlers.EmailHandlerBuild(config, emailService)

	log.Print("Setting up handlers")
	http.HandleFunc("/email", emailHandler.Handler)
	log.Print("Handlers successfully set up")

	log.Print("Starting email server")
	err = http.ListenAndServe(":8080", nil)
	log.Printf("Error occurred while running server: err=%+v", err)
}
