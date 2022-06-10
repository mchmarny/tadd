package todoist

import (
	"errors"
	"regexp"
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
			id, err := getItemID(apiLabelURL, apiToken, part[1:])
			if err != nil {
				return nil, err
			}
			t.Labels = append(t.Labels, *id)
			replace = append(replace, part)
		case "^":
			t.Due = &Due{}
			d := part[1:]
			replace = append(replace, part)
			pattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
			if pattern.MatchString(part[1:]) {
				t.Due.Date = &d
			} else {
				t.Due.String = &d
			}
		}
	}

	for _, r := range replace {
		content = strings.Replace(content, r, "", -1)
	}

	pattern := regexp.MustCompile(`\s+`)
	t.Content = pattern.ReplaceAllString(content, " ")

	return t, nil
}
