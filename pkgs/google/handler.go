package google

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Payload struct {
	Transition struct {
		WorkflowID     int    `json:"workflowId"`
		WorkflowName   string `json:"workflowName"`
		TransitionID   int    `json:"transitionId"`
		TransitionName string `json:"transitionName"`
		FromStatus     string `json:"from_status"`
		ToStatus       string `json:"to_status"`
	} `json:"transition"`
	Comment string `json:"comment"`
	User    struct {
		Self         string `json:"self"`
		Name         string `json:"name"`
		Key          string `json:"key"`
		EmailAddress string `json:"emailAddress"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
	} `json:"user"`
	Issue struct {
		ID     string `json:"id"`
		Self   string `json:"self"`
		Key    string `json:"key"`
		Fields struct {
			Parent struct {
				ID     string `json:"id"`
				Key    string `json:"key"`
				Self   string `json:"self"`
				Fields struct {
					Summary string `json:"summary"`
					Status  struct {
						Self           string `json:"self"`
						Description    string `json:"description"`
						IconURL        string `json:"iconUrl"`
						Name           string `json:"name"`
						ID             string `json:"id"`
						StatusCategory struct {
							Self      string `json:"self"`
							ID        int    `json:"id"`
							Key       string `json:"key"`
							ColorName string `json:"colorName"`
							Name      string `json:"name"`
						} `json:"statusCategory"`
					} `json:"status"`
					Priority struct {
						Self    string `json:"self"`
						IconURL string `json:"iconUrl"`
						Name    string `json:"name"`
						ID      string `json:"id"`
					} `json:"priority"`
					Issuetype struct {
						Self        string `json:"self"`
						ID          string `json:"id"`
						Description string `json:"description"`
						IconURL     string `json:"iconUrl"`
						Name        string `json:"name"`
						Subtask     bool   `json:"subtask"`
						AvatarID    int    `json:"avatarId"`
					} `json:"issuetype"`
				} `json:"fields"`
			} `json:"parent"`
			Assignee struct {
				Self         string `json:"self"`
				Name         string `json:"name"`
				Key          string `json:"key"`
				EmailAddress string `json:"emailAddress"`
				DisplayName  string `json:"displayName"`
				Active       bool   `json:"active"`
				TimeZone     string `json:"timeZone"`
			} `json:"assignee"`
			Status struct {
				Self           string `json:"self"`
				Description    string `json:"description"`
				IconURL        string `json:"iconUrl"`
				Name           string `json:"name"`
				ID             string `json:"id"`
				StatusCategory struct {
					Self      string `json:"self"`
					ID        int    `json:"id"`
					Key       string `json:"key"`
					ColorName string `json:"colorName"`
					Name      string `json:"name"`
				} `json:"statusCategory"`
			} `json:"status"`
			Archiveddate          interface{} `json:"archiveddate"`
			Customfield11017      interface{} `json:"customfield_11017"`
			Aggregatetimeestimate interface{} `json:"aggregatetimeestimate"`
			Customfield11090      interface{} `json:"customfield_11090"`
			Customfield11091      interface{} `json:"customfield_11091"`
			Creator               struct {
				Self         string `json:"self"`
				Name         string `json:"name"`
				Key          string `json:"key"`
				EmailAddress string `json:"emailAddress"`
				DisplayName  string `json:"displayName"`
				Active       bool   `json:"active"`
				TimeZone     string `json:"timeZone"`
			} `json:"creator"`
			Subtasks          []interface{} `json:"subtasks"`
			Reporter          interface{}   `json:"reporter"`
			Aggregateprogress struct {
				Progress int `json:"progress"`
				Total    int `json:"total"`
			} `json:"aggregateprogress"`
			Customfield11403 struct {
				Self     string `json:"self"`
				Value    string `json:"value"`
				ID       string `json:"id"`
				Disabled bool   `json:"disabled"`
				Child    struct {
					Self     string `json:"self"`
					Value    string `json:"value"`
					ID       string `json:"id"`
					Disabled bool   `json:"disabled"`
				} `json:"child"`
			} `json:"customfield_11403"`
			Issuetype struct {
				Self        string `json:"self"`
				ID          string `json:"id"`
				Description string `json:"description"`
				Name        string `json:"name"`
				Subtask     bool   `json:"subtask"`
				AvatarID    int    `json:"avatarId"`
			} `json:"issuetype"`
			Project struct {
				Self           string `json:"self"`
				ID             string `json:"id"`
				Key            string `json:"key"`
				Name           string `json:"name"`
				ProjectTypeKey string `json:"projectTypeKey"`
			} `json:"project"`
			Created          string      `json:"created"`
			Updated          string      `json:"updated"`
			Description      string      `json:"description"`
			Customfield10011 string      `json:"customfield_10011"`
			Customfield11500 string      `json:"customfield_11500"` //branch name
			Summary          string      `json:"summary"`
			Duedate          interface{} `json:"duedate"`
		} `json:"fields"`
	} `json:"issue"`
}

func Handler() {
	var Payload Payload
	b, e := os.ReadFile("Sample/task.json")
	if e != nil {
		print(e)
	}
	if err := json.Unmarshal(b, &Payload); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	fmt.Println(PrettyPrint(Payload.Issue.Fields.Issuetype.Name))

	issueType := Payload.Issue.Fields.Issuetype.Name
	issueKey := Payload.Issue.Key

	client := Init_client()
	switch issueType {
	case "task", "DevOps":
		err := client.GetMsg(issueKey)
		if err != nil {
			// in the create new msg with the same issuekey
			client.SendMsg(issueKey, issueKey)
		}
		log.Print("exist")
	case "subtask", "Sub-DevOps":
		parentIssue := Payload.Issue.Fields.Parent.Key

		err := client.GetMsg(issueKey)
		if err != nil {
			// in the create new msg with the same issuekey
			client.SendMsg(issueKey, parentIssue)
		}
		log.Print("exist")
	}

}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
