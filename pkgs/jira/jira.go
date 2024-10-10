package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	User struct {
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
	Changelog struct {
		ID    string `json:"id"`
		Items []struct {
			Field      string `json:"field"`
			Fieldtype  string `json:"fieldtype"`
			From       string `json:"from"`
			FromString string `json:"fromString"`
			To         string `json:"to"`
			ToString   string `json:"toString"`
		} `json:"items"`
	} `json:"changelog"`
}

func GetIssue(issueId string) (pl Payload, err error) {
	url := "https://task.sendo.vn/rest/api/2/issue/"+issueId
	method := "GET"

	jiraToken := os.Getenv("JIRA_TOKEN")

	client := &http.Client{}
	req, e := http.NewRequest(method, url, nil)

	if e != nil {
		fmt.Println(e)
		return pl, e
	}
	req.Header.Add("Authorization", "Bearer "+jiraToken)

	res, e := client.Do(req)
	if e != nil {
		fmt.Println(e)
		return pl, e
	}
	defer res.Body.Close()

	body, e := io.ReadAll(res.Body)
	if e != nil {
		fmt.Println(e)
		return pl, e
	}

	fmt.Println(string(body))
	if err := json.Unmarshal(body, &pl); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	return pl, nil
}
