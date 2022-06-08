package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	apiTaskURL         = "https://api.todoist.com/rest/v1/tasks"
	apiTokenEnvVarName = "TODOIST_API_TOKEN"
)

var (
	version  = "v0.0.1-default"
	commit   = ""
	date     = ""
	apiToken = ""

	task Task
)

type Task struct {
	ID           int64     `json:"id,omitempty"`
	ProjectID    int64     `json:"project_id,omitempty"`
	Content      string    `json:"content,omitempty"`
	Description  string    `json:"description,omitempty"`
	CommentCount int64     `json:"comment_count,omitempty"`
	Completed    bool      `json:"completed,omitempty"`
	Order        int64     `json:"order,omitempty"`
	Priority     int64     `json:"priority,omitempty"`
	Labels       []int     `json:"label_ids,omitempty"`
	SectionID    int64     `json:"section_id,omitempty"`
	ParentID     int64     `json:"parent_id,omitempty"`
	Creator      int64     `json:"creator,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	URL          string    `json:"url,omitempty"`
}

func addTask() (*Task, error) {
	b, err := json.Marshal(task)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal task")
	}

	req, err := http.NewRequest(http.MethodPost, apiTaskURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.Errorf("failed to add task: %s - %s", resp.Status, string(body))
	}

	var t Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, errors.Wrapf(err, "failed to decode response")
	}

	return &t, nil
}

func main() {
	flag.StringVar(&task.Content, "c", "", "The content of the task to create")
	flag.Int64Var(&task.ParentID, "p", 0, "ID of the project (optional, default: inbox)")
	flag.StringVar(&apiToken, "t", "", fmt.Sprintf("Todoist API token (default: $%s)", apiTokenEnvVarName))
	flag.Usage = func() {
		fmt.Printf("utility to quickly create Todoist task - %s (%s at %s)\n", version, commit, date)
		fmt.Println(" usage: ./td -c \"buy milk\"")
		flag.PrintDefaults()
	}
	flag.Parse()

	if apiToken == "" {
		apiToken = os.Getenv(apiTokenEnvVarName)
	}

	if task.Content == "" || apiToken == "" {
		fmt.Println("missing required arguments")
		fmt.Printf("content: %s, token:%s\n", task.Content, apiToken)
		flag.Usage()
		os.Exit(1)
	}

	t, err := addTask()
	if err != nil {
		fmt.Printf("error adding task: %s", err)
		os.Exit(1)
	}
	fmt.Printf("task created: %d\n", t.ID)
}
