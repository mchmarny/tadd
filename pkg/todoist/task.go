package todoist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiTaskCreateURL = "https://api.todoist.com/rest/v1/tasks"
)

type Task struct {
	Content      string     `json:"content,omitempty"` // required
	ID           *int64     `json:"id,omitempty"`
	ProjectID    *int64     `json:"project_id,omitempty"`
	Description  *string    `json:"description,omitempty"`
	CommentCount *int64     `json:"comment_count,omitempty"`
	Completed    *bool      `json:"completed,omitempty"`
	Order        *int64     `json:"order,omitempty"`
	Priority     *int       `json:"priority,omitempty"`
	Labels       []int64    `json:"label_ids,omitempty"`
	SectionID    *int64     `json:"section_id,omitempty"`
	ParentID     *int64     `json:"parent_id,omitempty"`
	Creator      *int64     `json:"creator,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
	URL          *string    `json:"url,omitempty"`
	DueDate      *string    `json:"due_date,omitempty"`
	DueString    *string    `json:"due_string,omitempty"`
}

func AddTask(apiToken, content string) (*Task, error) {
	task, err := parseTask(apiToken, content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task: %v", err)
	}

	resp, err := exec(http.MethodPost, apiTaskCreateURL, apiToken, task)
	if err != nil {
		return nil, fmt.Errorf("failed to post create task request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to add task: %s - %s", resp.Status, string(body))
	}

	var t Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, fmt.Errorf("unable to decode response: %v", err)
	}

	return &t, nil
}
