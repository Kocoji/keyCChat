package main

import (
	"github.com/joho/godotenv"
	"log"
	"notify-chat/cmd"
	// "notify-chat/pkgs/jira"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmd.Execute()
}
