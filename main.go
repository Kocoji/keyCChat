package main

import (
	"github.com/joho/godotenv"
	"log"
	"notify-chat/pkgs"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("KEYCLOAK_HOST")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	realm := os.Getenv("KC_REALM")
	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	kc, err := keycloak.InitKeyCloak(host,clientId,clientSecret,realm,user,password)
	kc.GetFUIdFromUId("kocojxi")
}
