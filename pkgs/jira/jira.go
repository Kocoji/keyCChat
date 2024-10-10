package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// the webhook payload
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
		Name         string `json:"name"`
		Key          string `json:"key"`
		EmailAddress string `json:"emailAddress"`
		DisplayName  string `json:"displayName"`
		Active       bool   `json:"active"`
	} `json:"user"`
	Issue struct {
		ID     string `json:"id"`
		Key    string `json:"key"`
		Fields struct {
			Parent struct {
				ID     string `json:"id"`
				Key    string `json:"key"`
				Fields struct {
					Summary string `json:"summary"`
					Status  struct {
						Description    string `json:"description"`
						IconURL        string `json:"iconUrl"`
						Name           string `json:"name"`
						ID             string `json:"id"`
						StatusCategory struct {
							ID        int    `json:"id"`
							Key       string `json:"key"`
							ColorName string `json:"colorName"`
							Name      string `json:"name"`
						} `json:"statusCategory"`
					} `json:"status"`
					Priority struct {
						IconURL string `json:"iconUrl"`
						Name    string `json:"name"`
						ID      string `json:"id"`
					} `json:"priority"`
					Issuetype struct {
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
				Name         string `json:"name"`
				Key          string `json:"key"`
				EmailAddress string `json:"emailAddress"`
				DisplayName  string `json:"displayName"`
				Active       bool   `json:"active"`
				TimeZone     string `json:"timeZone"`
			} `json:"assignee"`
			Status struct {
				Description    string `json:"description"`
				IconURL        string `json:"iconUrl"`
				Name           string `json:"name"`
				ID             string `json:"id"`
				StatusCategory struct {
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
				Value    string `json:"value"`
				ID       string `json:"id"`
				Disabled bool   `json:"disabled"`
				Child    struct {
					Value    string `json:"value"`
					ID       string `json:"id"`
					Disabled bool   `json:"disabled"`
				} `json:"child"`
			} `json:"customfield_11403"`
			Issuetype struct {
				ID          string `json:"id"`
				Description string `json:"description"`
				Name        string `json:"name"`
				Subtask     bool   `json:"subtask"`
				AvatarID    int    `json:"avatarId"`
			} `json:"issuetype"`
			Project struct {
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

type Response struct {
	Expand string `json:"expand"`
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Summary          string `json:"summary"`
		Customfield11403 struct {
			Value    string `json:"value"`
			ID       string `json:"id"`
			Disabled bool   `json:"disabled"`
			Child    struct {
				Value    string `json:"value"`
				ID       string `json:"id"`
				Disabled bool   `json:"disabled"`
			} `json:"child"`
		} `json:"customfield_11403"`
		Creator struct {
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			DisplayName  string `json:"displayName"`
			Active       bool   `json:"active"`
			TimeZone     string `json:"timeZone"`
		} `json:"creator"`
		Description string `json:"description"`
		Assignee    struct {
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			DisplayName  string `json:"displayName"`
			Active       bool   `json:"active"`
			TimeZone     string `json:"timeZone"`
		} `json:"assignee"`
		Reporter struct {
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			DisplayName  string `json:"displayName"`
			Active       bool   `json:"active"`
			TimeZone     string `json:"timeZone"`
		} `json:"reporter"`
		Status struct {
			Description    string `json:"description"`
			IconURL        string `json:"iconUrl"`
			Name           string `json:"name"`
			ID             string `json:"id"`
			StatusCategory struct {
				ID        int    `json:"id"`
				Key       string `json:"key"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
	} `json:"fields"`
}

func GetIssue(issueId string) (res Response, err error) {

	extraFields := "?fields=status,creator,assignee,customfield_11403,summary,description"
	jiraToken := os.Getenv("JIRA_TOKEN")

	url := "https://task.sendo.vn/rest/api/2/issue/" + issueId + extraFields
	method := "GET"

	client := &http.Client{}
	req, e := http.NewRequest(method, url, nil)

	if e != nil {
		fmt.Println(e)
		return res, e
	}
	req.Header.Add("Authorization", "Bearer "+jiraToken)

	response, e := client.Do(req)
	if e != nil {
		fmt.Println(e)
		return res, e
	}
	defer response.Body.Close()

	body, e := io.ReadAll(response.Body)
	if e != nil {
		fmt.Println(e)
		return res, e
	}

	// fmt.Println(string(body))
	if err := json.Unmarshal(body, &res); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	return res, nil
}
