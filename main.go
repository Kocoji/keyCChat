package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"github.com/joho/godotenv"
)


type ClientParams struct {
	client *gocloak.GoCloak
	ctx    context.Context
	token  *gocloak.JWT
	realm  string
}

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

	client := gocloak.NewClient(host)
	ctx := context.Background()
	token, err := client.Login(ctx, clientId, clientSecret, realm, user, password)
	if err != nil {
		log.Println(err)
		panic("Something wrong with the credentials or url")
	}

	params := ClientParams{
		client: client,
		ctx:    ctx,
		token:  token,
		realm:  realm,
	}
	uid, err := getUserId(params, "kocoji")
	if err != nil {
		log.Println("It's seem like the user is not exist, please try another name")
	}
	userFedId,err := getFedUserId(params, uid)
	if err != nil {
		log.Println("Not exist ", err)
	}

	log.Println("Fed userId:", userFedId)

}

// this func get the list of
func getUserId(client ClientParams, username string) (string, error) {
	user, err := client.client.GetUsers(client.ctx, client.token.AccessToken, client.realm, gocloak.GetUsersParams{Username: gocloak.StringP(username)})
	if err != nil {
		log.Printf("It's seem like the %s is not exist!", username)
	}
	if len(user) > 0 {
		return *user[0].ID, nil
	} else {
		return "", errors.New("empty name")
	}
}

func getFedUserId(client ClientParams, userId string) (string, error) {
	userFedId, err := client.client.GetUserFederatedIdentities(client.ctx, client.token.AccessToken, client.realm, userId)
	if err != nil {
		log.Println("It's seem like the federal user is not exist! May be this user didn't connect into any Id Provider")
	}
	if len(userFedId) > 0 {
		return *userFedId[0].UserID, nil
	} else {
		return "", errors.New("empty name")
	}
}
