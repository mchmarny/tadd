package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mchmarny/tadd/pkg/todoist"
)

const (
	envVarName = "TODOIST_API_TOKEN" // because the word "token" in variable name makes linter unhappy
)

var (
	version = "v0.0.1-default"

	commit   string
	date     string
	apiToken string
	content  string
)

func main() {
	flag.StringVar(&content, "c", "", "The content of the task to create")
	flag.StringVar(&apiToken, "t", "", fmt.Sprintf("Todoist API token (default: $%s)", envVarName))
	flag.Usage = func() {
		fmt.Printf("utility to quickly create Todoist task - %s (%s at %s)\n", version, commit, date)
		fmt.Println(" usage: tadd -c \"buy milk\"")
		flag.PrintDefaults()
	}
	flag.Parse()

	if apiToken == "" {
		apiToken = os.Getenv(envVarName)
	}

	if content == "" || apiToken == "" {
		fmt.Println("missing required arguments")
		flag.Usage()
		os.Exit(1)
	}

	t, err := todoist.AddTask(apiToken, content)
	if err != nil {
		fmt.Printf("error adding task: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Task added: %s\n", *t.URL)
}
