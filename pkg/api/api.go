package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go_jira_logger/pkg/config"
	assert "go_jira_logger/pkg/utils"

	"github.com/spf13/viper"
)

type JiraIssueType struct {
	Self           string `json:"self"`
	ID             string `json:"id"`
	Description    string `json:"description"`
	IconURL        string `json:"iconUrl"`
	Name           string `json:"name"`
	Subtask        bool   `json:"subtask"`
	AvatarID       int    `json:"avatarId"`
	HierarchyLevel int    `json:"hierarchyLevel"`
}

type JiraStatusCategory struct {
	Self      string `json:"self"`
	ID        int    `json:"id"`
	Key       string `json:"key"`
	ColorName string `json:"colorName"`
	Name      string `json:"name"`
}

type JiraStatus struct {
	Self           string             `json:"self"`
	Description    string             `json:"description"`
	IconURL        string             `json:"iconUrl"`
	Name           string             `json:"name"`
	ID             string             `json:"id"`
	StatusCategory JiraStatusCategory `json:"statusCategory"`
}
type JiraIssueFields struct {
	Summary     string        `json:"summary"`
	Issuetype   JiraIssueType `json:"issuetype"`
	Description interface{}   `json:"description"`
	Status      JiraStatus    `json:"status"`
}

type JiraIssue struct {
	Expand string          `json:"expand"`
	ID     string          `json:"id"`
	Self   string          `json:"self"`
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

type IssuesResponse struct {
	Expand     string      `json:"expand"`
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []JiraIssue `json:"issues"`
}

func FetchIssues() IssuesResponse {
	config.ReadConfigFile()

	baseUrl := viper.Get("jira.board_url")
	searchEnpoint := "/rest/api/3/search"

	u, err := url.Parse(fmt.Sprintf("%s%s", baseUrl, searchEnpoint))

	assert.NoError(err, "Error parsing URL")

	jiraEmail := viper.Get("jira.email")
	jiraToken := viper.Get("jira.token")

	jql := `sprint in openSprints() AND (status = "In Progress" OR status="Peer Review") AND (issuetype = Task OR issuetype = Bug) AND assignee = currentUser() ORDER BY created DESC`
	fields := "summary,description,status,issuetype"

	params := url.Values{}
	params.Add("jql", jql)
	params.Add("fields", fields)

	u.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	assert.NoError(err, "Could not construct http request")

	auth := fmt.Sprintf("%s:%s", jiraEmail, jiraToken)

	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

	client := &http.Client{}

	response, err := client.Do(req)

	assert.NoError(err, "Error fetching the jira issues")

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	assert.NoError(err, "Reponse body could not be parsed")

	var issueResponse IssuesResponse
	err = json.Unmarshal(body, &issueResponse)

	assert.NoError(err, "Could not parse json body")

	return issueResponse
}
