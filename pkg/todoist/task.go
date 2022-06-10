package todoist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiTaskURL = "https://api.todoist.com/rest/v1/tasks"
)

type Task struct {
	Content      string     `json:"content,omitempty"` // required
	ID           *int64     `json:"id,omitempty"`
	ProjectID    *int64     `json:"project_id,omitempty"`
	Description  *string    `json:"description,omitempty"`
	CommentCount *int64     `json:"comment_count,omitempty"`
	Completed    *bool      `json:"completed,omitempty"`
	Order        *int64     `json:"order,omitempty"`
	Priority     *int64     `json:"priority,omitempty"`
	Labels       []int      `json:"label_ids,omitempty"`
	SectionID    *int64     `json:"section_id,omitempty"`
	ParentID     *int64     `json:"parent_id,omitempty"`
	Creator      *int64     `json:"creator,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
	URL          *string    `json:"url,omitempty"`
}

func AddTask(apiToken, content string) (*Task, error) {
	if apiToken == "" {
		return nil, errors.New("missing required argument: apiToken")
	}

	if content == "" {
		return nil, errors.New("missing required argument: content")
	}

	task := &Task{
		Content: content,
	}

	b, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("error marshalling task: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiTaskURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error posting request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to add task: %s - %s", resp.Status, string(body))
	}

	var t Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	return &t, nil
}
