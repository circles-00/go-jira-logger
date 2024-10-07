package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"go_jira_logger/pkg/config"
	assert "go_jira_logger/pkg/utils/assert"
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

type WorklogPayload struct {
	TimeSpent string `json:"timeSpent"`
	Started   string `json:"started"`
}

type AddWorklogResponse struct {
	Id string `json:"id"`
}

type AttachTagsToWorklogPayload struct {
	Tags []string `json:"tags"`
}

func FetchIssues() IssuesResponse {
	config.ReadConfigFile()

	baseUrl := config.GetBoardUrl()
	searchEnpoint := "/rest/api/3/search"

	u, err := url.Parse(fmt.Sprintf("%s%s", baseUrl, searchEnpoint))

	assert.NoError(err, "Error parsing URL")

	jiraEmail := config.GetJiraEmail()
	jiraToken := config.GetJiraToken()

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

func SubmitWorklog(issueKey string, tags []string, worklogPayload WorklogPayload) AddWorklogResponse {
	config.ReadConfigFile()

	baseUrl := config.GetBoardUrl()
	worklogEnpoint := fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueKey)

	u, err := url.Parse(fmt.Sprintf("%s%s", baseUrl, worklogEnpoint))

	assert.NoError(err, "Error parsing URL")

	jiraEmail := config.GetJiraEmail()
	jiraToken := config.GetJiraToken()

	params := url.Values{}

	u.RawQuery = params.Encode()

	requestBody, err := json.Marshal(worklogPayload)

	assert.NoError(err, "Error serializing request body")

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBuffer(requestBody))

	assert.NoError(err, "Could not construct http request")

	auth := fmt.Sprintf("%s:%s", jiraEmail, jiraToken)

	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

	client := &http.Client{}

	response, err := client.Do(req)

	assert.NoError(err, "Error adding worklog")

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	assert.NoError(err, "Reponse body could not be parsed")

	var addWorklogResponse AddWorklogResponse
	err = json.Unmarshal(body, &addWorklogResponse)

	assert.NoError(err, "Could not parse json body")

	AttachTagsToWorklog(issueKey, addWorklogResponse.Id, tags)

	return addWorklogResponse
}

func AttachTagsToWorklog(issueKey string, worklogId string, tags []string) {
	config.ReadConfigFile()

	baseUrl := config.GetBoardUrl()
	worklogEnpoint := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s/properties/com.gebsun.plugins.work.tags", issueKey, worklogId)

	u, err := url.Parse(fmt.Sprintf("%s%s", baseUrl, worklogEnpoint))

	assert.NoError(err, "Error parsing URL")

	jiraEmail := config.GetJiraEmail()
	jiraToken := config.GetJiraToken()

	params := url.Values{}

	u.RawQuery = params.Encode()

	requestBody, err := json.Marshal(AttachTagsToWorklogPayload{
		Tags: tags,
	})

	assert.NoError(err, "Error serializing request body")

	req, err := http.NewRequest(http.MethodPut, u.String(), bytes.NewBuffer(requestBody))

	assert.NoError(err, "Could not construct http request")

	auth := fmt.Sprintf("%s:%s", jiraEmail, jiraToken)

	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

	client := &http.Client{}

	response, err := client.Do(req)

	assert.NoError(err, "Error attaching worklog tags")

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	assert.NoError(err, "Reponse body could not be parsed")

	var addTagsToWorklogResponse any
	err = json.Unmarshal(body, &addTagsToWorklogResponse)

	assert.NoError(err, "Could not parse json body")

	log.Printf("Successfully attached tags to worklog with task: %s and id %s", issueKey, worklogId)
}
