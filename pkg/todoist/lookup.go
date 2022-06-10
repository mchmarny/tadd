package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiProjectURL = "https://api.todoist.com/rest/v1/projects"
	apiLabelURL   = "https://api.todoist.com/rest/v1/labels"
)

type Item struct {
	ID   *int64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func getItemID(url, token, name string) (*int64, error) {
	if url == "" {
		return nil, errors.New("missing required argument: url")
	}

	if token == "" {
		return nil, errors.New("missing required argument: token")
	}

	if name == "" {
		return nil, errors.New("missing required argument: projectName")
	}

	resp, err := exec(http.MethodGet, url, token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to post get projects request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get projects: %s - %s", resp.Status, string(body))
	}

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("unable to decode response: %v", err)
	}

	for _, item := range items {
		if strings.EqualFold(item.Name, name) {
			return item.ID, nil
		}
	}

	return nil, nil
}

func createItem(url, token, name string) (*int64, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}

	if token == "" {
		return nil, errors.New("missing required argument: token")
	}

	if name == "" {
		return nil, errors.New("missing required argument: projectName")
	}

	resp, err := exec(http.MethodPost, url, token, &Item{Name: name})
	if err != nil {
		return nil, fmt.Errorf("error posting get projects request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get projects: %s - %s", resp.Status, string(body))
	}

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return item.ID, nil
}
