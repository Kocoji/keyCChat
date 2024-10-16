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
	IssueId     string
	ParentId    string
	Descript    string
	Summary     string
	Psummary    string
	ChangelogId string
	UserFedId   string
	Status      string
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
func (c *Client) SendMsg(j Msg, parent bool) error {

	threadKey := ""
	spacepath := "spaces/" + c.space
	jira_uri := os.Getenv("JIRA_HOST") + "/browse/" + j.IssueId
	sumary := ""
	msgId := "client-" + strings.ToLower(j.IssueId)
	if parent || (j.ParentId != "") {
		// msgId = "client-" + strings.ToLower(j.IssueId) + "-" + j.ChangelogId
		threadKey = j.ParentId
		sumary = j.Psummary
	} else {
		threadKey = j.IssueId
		sumary = j.Summary
		// msgId = "client-" + strings.ToLower(j.IssueId)
	}
	// if parrent not empty, send a message to Parent's thread message.
	// if j.ParentId != "" {
	// 	// uncommend in
	// 	// msgId = "client-" + strings.ToLower(j.IssueId) + "-" + j.ChangelogId
	// 	threadKey = j.ParentId
	// }
	fmt.Println("msg id ", msgId)
	msg := chat.Message{
		Text: "Hi <users/" + j.UserFedId + "> This task need your action \nCurent Status:" + j.Status + "\n" + sumary + "\nLink:" + jira_uri,
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

func (c *Client) UpdateMsg(j Msg, thread bool) {
	space := c.space
	threadKey := j.IssueId
	msgId := ""
	if thread {
		msgId = "client-" + strings.ToLower(j.IssueId) + "-" + j.ChangelogId
		threadKey = j.ParentId
	} else {
		msgId = "client-" + strings.ToLower(j.IssueId)
	}
	name := "spaces/" + space + "/messages/" + msgId
	jira_uri := os.Getenv("JIRA_HOST") + "/browse/" + j.IssueId

	msg := chat.Message{
		Text: "Hi <users/" + j.UserFedId + "> This task need your action \nCurent Status:" + j.Status + "\n" + j.Psummary + "\nLink:" + jira_uri,
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
	return res, true
}

func (c *Client) DelMsg(msgId string) error {
	_, e := c.client.Spaces.Messages.Delete("client-" + strings.ToLower(msgId)).Do()
	if e != nil {
		log.Println(e)
		return fmt.Errorf("%v", e)
	}
	return nil
}
