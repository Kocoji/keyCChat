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
	space  string
}

type Msg struct {
	msg chat.Message
}
func Init_client() Client {
	space := os.Getenv("GC_SPACE")

	ctx := context.Background()
	client, err := chat.NewService(ctx, option.WithCredentialsFile("credentials.json"), option.WithScopes(chat.ChatBotScope))
	if err != nil {
		fmt.Println("Oops", err)
	}
	// client.Spaces.Messages
	return Client{
		client: client,
		space:  space,
	}
}

// this func use to send a new message, so issueType: Task, DevOps use this func.
func (c *Client) SendMsg(issueId string, parentIss string, mesg string, changelogId string) error {
	msgId := "client-" + strings.ToLower(issueId)

	threadKey := issueId
	spacepath := "spaces/" + c.space

	if parentIss != "" {
		msgId = "client-" + strings.ToLower(issueId) + "-" + changelogId
		threadKey = parentIss
	}
	fmt.Println("msg id ",msgId)
	msg := chat.Message{
		Text: mesg,
		Thread: &chat.Thread{
			ThreadKey: threadKey,
		},
	}
	_, e := c.client.Spaces.Messages.Create(spacepath, &msg).MessageReplyOption("REPLY_MESSAGE_FALLBACK_TO_NEW_THREAD").MessageId(msgId).Do()
	if e != nil {
		log.Fatal(e)
	}
	return nil
}

func (c *Client) UpdateMsg(threadKey string, msid string) {
	space := c.space
	msgId := "client-" + strings.ToLower(msid)
	name := "spaces/" + space + "/messages/" + msgId

	msg := chat.Message{
		Text: "Updated",
		Thread: &chat.Thread{
			ThreadKey: threadKey,
		},
	}
	r, err := c.client.Spaces.Messages.Update(name, &msg).UpdateMask("text").Do()
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("res: ", r)
}

// Get message by messageId and return data & bool true (exist)/false(not exist)
func (c *Client) GetMsg(issueId string) (*chat.Message, bool) {
	msgId := "client-" + strings.ToLower(issueId)
	space := c.space
	name := "spaces/" + space + "/messages/" + msgId
	res, err := c.client.Spaces.Messages.Get(name).Do()
	if err != nil {
		log.Println(err)
		return nil, false
	}
	// fmt.Print(res)
	return res, true
}
