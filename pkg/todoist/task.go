package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
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
