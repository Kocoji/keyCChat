package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"notify-chat/pkgs/google"
	"notify-chat/pkgs/jira"
	"notify-chat/pkgs/keycloak"
	"os"
)

// I need to put link to docs here.
func Handler() error {
	Payload := jira.Payload{}
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
	userEmail := Payload.Issue.Fields.Assignee.EmailAddress

	client := google.Init_client()

	kc, err := keycloak.InitKeyCloak()
	if err != nil {
		log.Fatalln("Problem when try to access Keycloak")
	}
	fedUserId:= kc.GetFUIdFromUId(userEmail)
	kc.Logout()

	jiraData := google.Msg{
		IssueId:     issueKey,
		Summary:     issueName,
		ChangelogId: changeLog,
		Descript: Payload.Issue.Fields.Description,
		UserFedId: fedUserId,
		Status: Payload.Issue.Fields.Status.Name,
	}

	switch issueType {
	case "Task", "DevOps":
		_, exist := client.GetMsg(issueKey)
		if exist {
			log.Println("exist")
			// Update the main message and send a new one in the thread, mentioning the assignee
			client.UpdateMsg(jiraData)
			client.SendMsg(jiraData, true)
		} else {
			// in the create new msg with the same issuekey
			client.SendMsg(jiraData, false)
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
			log.Println("This exist")

			// Update Parent task's status
			client.UpdateMsg(jiraData)
			// then send subtask as thread message.
			client.SendMsg(jiraData, false)
		} else {
			jiraData.Summary = parentIssueSum
			// the create new msg with the same issuekey
			// Send Parent Task message
			client.SendMsg(jiraData, false)
			// then, send subtask as thread message
			client.SendMsg(jiraData, true)
		}
	}
	return nil
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
