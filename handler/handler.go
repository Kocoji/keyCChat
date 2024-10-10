package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"notify-chat/pkgs/jira"
	"notify-chat/pkgs/google"
)


// I need to put link to docs here.
func Handler() error {
	var Payload jira.Payload
	b, e := os.ReadFile("Sample/task.json")
	if e != nil {
		print(e)
	}
	if err := json.Unmarshal(b, &Payload); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	fmt.Println(PrettyPrint(Payload))

	issueType := Payload.Issue.Fields.Issuetype.Name
	issueKey := Payload.Issue.Key

	issueName := Payload.Issue.Fields.Summary
	changeLog := Payload.Changelog.ID

	client := google.Init_client()
	switch issueType {
	case "Task", "DevOps":
		err := client.GetMsg(issueKey)
		if err != nil {
			// in the create new msg with the same issuekey
			client.SendMsg(issueKey,"", issueName,"")
		} else {
			log.Print("exist")
			// update main message and send a new message in Thread that mention assignee use.
			client.UpdateMsg(issueKey, issueKey)
			client.SendMsg(issueKey, "", issueName,changeLog)
		}
	case "Subtask", "Sub-DevOps":
		parentIssue := Payload.Issue.Fields.Parent.Key
		err := client.GetMsg(parentIssue)
		if err != nil {
			res, e := jira.GetIssue(parentIssue)
			if e !=nil {
				log.Println(e)
			}

			parentChangeLog := res.Changelog.ID
			// in the create new msg with the same issuekey
			client.SendMsg(parentIssue, parentIssue, issueName,parentChangeLog)
		} else {
			log.Print("exist")
			client.UpdateMsg(issueKey, issueKey)
			client.SendMsg(issueKey, parentIssue, issueName,changeLog)
		}
	}
	return nil
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
