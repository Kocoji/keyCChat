package keycloak

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Nerzal/gocloak/v13"
)

type KeyCloakClient struct {
	client *gocloak.GoCloak
	ctx    context.Context
	token  *gocloak.JWT
	realm  string
}

type User struct {
	UserId   string
	Username string
}

func InitKeyCloak(host, clientId, clientSecret, realm, user, password string) (KeyCloakClient, error) {
	client := gocloak.NewClient(host)
	ctx := context.Background()
	token, err := client.Login(ctx, clientId, clientSecret, realm, user, password)
	if err != nil {
		log.Fatal("Cannot init keycloak client")
		panic(err)
	}
	return KeyCloakClient{
		client: client,
		ctx:    ctx,
		token:  token,
		realm:  realm}, nil
}

func (c *KeyCloakClient) GetFUIdFromUId(username string) string {

	uid, err := c.getUserId(username)
	if err != nil {
		log.Println(err)
		return ""
	}
	userFedId, err := c.getFedUserId(uid)
	if err != nil {
		log.Println("Not exist ", err)
	}
	return userFedId
}

// this func get the list of
func (c *KeyCloakClient) getUserId(username string) (string, error) {
	user, err := c.client.GetUsers(c.ctx, c.token.AccessToken, c.realm, gocloak.GetUsersParams{Username: gocloak.StringP(username)})
	if err != nil {
		fmt.Println("client error:",err)
		return "",fmt.Errorf("%v",err)
	}
	if len(user) == 0 {
		return "", errors.New("it seems like the user does not exist, please try another username")
	} else {
		return *user[0].ID, nil
	}
}

func (c *KeyCloakClient) getFedUserId( userId string) (string, error) {
	userFedId, err := c.client.GetUserFederatedIdentities(c.ctx, c.token.AccessToken, c.realm, userId)
	if err != nil {
		log.Println("client error:",err)
		return "",fmt.Errorf("%v",err)
	}
	if len(userFedId) == 0 {
		return "", errors.New("oops! it seems this user didn't connect to any id provider")

	} else {
		return *userFedId[0].UserID, nil
	}
}
