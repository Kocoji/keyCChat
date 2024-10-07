package google

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

type Client struct {
	client *chat.Service
	space string
}

type Fields struct {
	spaces string;
}

func Init_client() Client {
	space := os.Getenv("GC_SPACE")

	ctx := context.Background()
	client, err := chat.NewService(ctx, option.WithCredentialsFile("credentials.json"), option.WithScopes(chat.ChatBotScope))
	if err != nil {
		fmt.Println("OOPS", err)
	}
	fmt.Printf("%+v\n", client)

	// client.Spaces.Messages
	return Client{
		client: client,
		space: space,
	}
}

// this func use to send a new message, so issueType: Task, DevOps use this func.
func (c *Client) SendMsg(threadKey string, msid string) error {
	spacepath := "spaces/" + c.space
	msgId := "client-" + strings.ToLower(msid)
	msg := chat.Message{
		Text: "Hello world!6733",
		Thread: &chat.Thread{
			ThreadKey: threadKey,
		},
	}
	_, e := c.client.Spaces.Messages.Create(spacepath, &msg).MessageReplyOption("REPLY_MESSAGE_FALLBACK_TO_NEW_THREAD").MessageId(msgId).Do()
	if e != nil {
		log.Fatal(e)
	}
	return nil
	// fmt.Println("res: ", r)
}

func (c *Client) UpdateMsg() {
	space := c.space
	msgId := "client-6732"
	name := "spaces/" + space + "/messages/" + msgId

	msg := chat.Message{
		Text: "Hello ahihi  d!s",
		Thread: &chat.Thread{
			ThreadKey: "ahihi123",
		},
	}
	r, err := c.client.Spaces.Messages.Update(name, &msg).UpdateMask("text").Do()
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("res: ", r)
}

func (c *Client) GetMsg(msid string) error {
	msgId := "client-" + strings.ToLower(msid)
	space := c.space
	name := "spaces/" + space + "/messages/" + msgId
	res,err := c.client.Spaces.Messages.Get(name).Do()
	if err !=nil {
		fmt.Print(err)
		return fmt.Errorf("problem")
	}
	fmt.Print(res)
	return nil
}
