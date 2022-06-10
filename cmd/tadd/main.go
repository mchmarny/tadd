package main

import (
	"flag"
	"fmt"
	"os"
	"time"

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
		fmt.Println()
		fmt.Printf("tadd - Utility to quickly create Todoist task - %s (%s at %s)\n", version, commit, date)
		fmt.Println()
		fmt.Println(" Usage: tadd -c \"buy milk ^monday #personal @shopping *4\"")
		fmt.Println(" Result: task 'buy milk' due on 'monday' with a label 'shopping' in project 'personal' with highest priority")
		fmt.Println()
		fmt.Println(" Note: ")
		fmt.Println("    projects default to 'inbox' if don't exist or not specified")
		fmt.Println("    labels (prefix:@) will be created if don't exist")
		fmt.Println("    priority (prefix:*) span from 1-normal to 4-high")
		fmt.Printf("    due dates (prefix:^) can be relative (e.g. ^tomorrow) or absolute (e.g. ^%s)\n", time.Now().Format("2006-01-02"))
		fmt.Println()
		fmt.Println(" Arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if apiToken == "" {
		apiToken = os.Getenv(envVarName)
	}

	if apiToken == "" || content == "" {
		fmt.Println("missing arguments")
		flag.Usage()
		os.Exit(1)
	}

	t, err := todoist.AddTask(apiToken, content)
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	fmt.Println(*t.URL)
}
