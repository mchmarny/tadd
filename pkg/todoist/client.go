package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func exec(method, url, apiToken string, content interface{}) (*http.Response, error) {
	var body io.Reader

	if content != nil {
		b, err := json.Marshal(content)
		if err != nil {
			return nil, fmt.Errorf("error marshaling content: %v", err)
		}
		body = bytes.NewBuffer(b)

		// fmt.Printf("%s %s\n%s\n", method, url, string(b))
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	c := &http.Client{}
	return c.Do(req)
}
