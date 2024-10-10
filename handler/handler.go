package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"notify-chat/pkgs/google"
	"notify-chat/pkgs/jira"
	"os"
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
		_, exist := client.GetMsg(issueKey)
		if exist {
			log.Println("exist")
			// update main message and send a new message in Thread that mention assignee use.
			client.UpdateMsg(issueKey, issueKey)
			client.SendMsg(issueKey, "", issueName, changeLog)
		} else {
			// in the create new msg with the same issuekey
			client.SendMsg(issueKey, "", issueName, "")
		}
	case "Subtask", "Sub-DevOps":
		parentIssue := Payload.Issue.Fields.Parent.Key
		res, e := jira.GetIssue(parentIssue)
		parentIssueSum := res.Fields.Summary

		if e != nil {
			log.Println(e)
		}
		_, exist := client.GetMsg(parentIssue)
		if exist {
			log.Println("exist")

			client.UpdateMsg(parentIssue, parentIssue)
			client.SendMsg(issueKey, parentIssue, issueName, changeLog)

		} else {
			// in the create new msg with the same issuekey
			client.SendMsg(parentIssue, "", parentIssueSum, "")
		}
	}
	return nil
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
