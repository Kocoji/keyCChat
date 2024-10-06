package google

import (
	"context"
	"fmt"

	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

type Client struct {
	client *chat.Service
}

func Init_client() Client {
	ctx := context.Background()
	client, err := chat.NewService(ctx, option.WithCredentialsFile("credentials.json"), option.WithScopes(chat.ChatBotScope))
	if err != nil {
		fmt.Println("OOPS", err)
	}
	fmt.Printf("%+v\n", client)

	// client.Spaces.Messages
	return Client{
		client: client,
	}
}

func (c *Client) SendMsg(){
	space := "spaces/"+"AAAAqZC-F5E"

	msg := chat.Message{
		Text: "Hello world!123",
		Thread: &chat.Thread{
			ThreadKey: "ahihi123",
		},
		ThreadReply: true,
	}
	request, err := c.client.Spaces.Messages.Create(space, &msg).MessageId("client-ab2cd1234").Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res: ",request)
}


func (c *Client) UpdateMsg(){
	space := "AAAAqZC-F5E"
	msgId := "client-ab2cd1234"
	name := "spaces/"+space+ "/messages/"+msgId

	msg := chat.Message{
		Text: "Hello ahihi  d!s",
		Thread: &chat.Thread{
			ThreadKey: "ahihi123",
		},
		ThreadReply: true,
	}
	r, err := c.client.Spaces.Messages.Update(name, &msg).UpdateMask("text").Do()
	if err != nil {
		fmt.Println("err:",err)
	}
	fmt.Println("res: ",r)
}


func (c *Client) SendThreadMsg(){
	space := "spaces/"+"AAAAqZC-F5E"
	msg := chat.Message{
		Text: "Hello world!123",
		Thread: &chat.Thread{
			ThreadKey: "ahihi123",
		},
		ThreadReply: true,
	}
	r, err := c.client.Spaces.Messages.Create(space, &msg).MessageReplyOption("REPLY_MESSAGE_FALLBACK_TO_NEW_THREAD").Do()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res: ",r)
}
