package keycloak

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type KeyCloakClient struct {
	client   *gocloak.GoCloak
	ctx      context.Context
	token    *gocloak.JWT
	realm    string
	clientId string
	clientSecret string
}

type User struct {
	UserId   string
	Username string
}

func InitKeyCloak() (KeyCloakClient, error) {
	host := os.Getenv("KEYCLOAK_HOST")
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	realm := os.Getenv("KC_REALM")
	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	client := gocloak.NewClient(host)
	ctx := context.Background()
	token, err := client.Login(ctx, clientId, clientSecret, "master", user, password)
	if err != nil {
		log.Fatal("Cannot init keycloak client")
		panic(err)
	}
	return KeyCloakClient{
		client:   client,
		ctx:      ctx,
		token:    token,
		realm:    realm,
		clientId: clientId,
		clientSecret: clientSecret,}, nil
}

func (k *KeyCloakClient) GetFUIdFromUId(username string) string {
	uid, err := k.getUserId(username)
	if err != nil {
		log.Println(err)
		return ""
	}
	userFedId, err := k.getFedUserId(uid)
	if err != nil {
		log.Println("Not exist ", err)
	}
	return userFedId
}

// this func get the list of
func (k *KeyCloakClient) getUserId(username string) (string, error) {
	user, err := k.client.GetUsers(k.ctx, k.token.AccessToken, k.realm, gocloak.GetUsersParams{Username: gocloak.StringP(username)})
	if err != nil {
		fmt.Println("client error:", err)
		return "", fmt.Errorf("%v", err)
	}
	if len(user) == 0 {
		return "", errors.New("it seems like the user does not exist, please try another username")
	} else {
		return *user[0].ID, nil
	}
}

func (k *KeyCloakClient) Logout() error {
	e := k.client.Logout(k.ctx,k.clientId, k.clientSecret,"master",k.token.RefreshToken )
	if e !=nil {
		println(e)
	}
	return nil
}

func (k *KeyCloakClient) getFedUserId(userId string) (string, error) {
	userFedId, err := k.client.GetUserFederatedIdentities(k.ctx, k.token.AccessToken, k.realm, userId)
	if err != nil {
		log.Println("client error:", err)
		return "", fmt.Errorf("%v", err)
	}
	if len(userFedId) == 0 {
		return "", errors.New("oops! it seems this user didn't connect to any id provider")

	} else {
		return *userFedId[0].UserID, nil
	}
}
