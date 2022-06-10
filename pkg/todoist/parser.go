package todoist

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parseTask(apiToken, content string) (*Task, error) {
	if content == "" {
		return nil, errors.New("missing required argument: content")
	}

	if apiToken == "" {
		return nil, errors.New("missing required argument: apiToken")
	}

	t := &Task{}

	parts := strings.Split(content, " ")
	replace := make([]string, 0)

	for _, part := range parts {
		if len(part) <= 1 {
			continue
		}

		prefix := part[:1]
		switch prefix {
		case "#":
			id, err := getItemID(apiProjectURL, apiToken, part[1:])
			if err != nil {
				return nil, err
			}
			t.ProjectID = id
			replace = append(replace, part)
		case "@":
			// get existing label id
			id, err := getItemID(apiLabelURL, apiToken, part[1:])
			if err != nil {
				return nil, err
			}

			// if label doesn't exist, create it
			if id == nil {
				id, err = createItem(apiLabelURL, apiToken, part[1:])
				if err != nil {
					return nil, err
				}
			}

			t.Labels = append(t.Labels, *id)
			replace = append(replace, part)
		case "^":
			d := part[1:]
			replace = append(replace, part)
			pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
			if pattern.MatchString(part[1:]) {
				t.DueDate = &d
			} else {
				t.DueString = &d
			}
		case "*":
			i, err := strconv.Atoi(part[1:])
			if err != nil {
				return nil, fmt.Errorf("invalid priority: %s", part[1:])
			}
			if i < 1 || i > 4 {
				return nil, fmt.Errorf("priority must be 1-4, got: %d", i)
			}
			replace = append(replace, part)
			t.Priority = &i
		}
	}

	for _, r := range replace {
		content = strings.Replace(content, r, "", -1)
	}

	pattern := regexp.MustCompile(`\s+`)
	t.Content = pattern.ReplaceAllString(content, " ")

	return t, nil
}
