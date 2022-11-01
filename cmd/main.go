package main

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"github.com/opsgenie/opsgenie-go-sdk-v2/incident"
	"log"
	"os"
)

func main() {
	//Extract all arguments
	allArguments := os.Args[1:]

	//Did we get the correct number of arguments
	if len(allArguments) != 5 {
		log.Fatalf("wrong number of arguments passed, got %d arguments.", len(allArguments))
	}

	//Validate the priority string
	var priority = allArguments[1]
	if len(priority) != 2 {
		log.Fatal("priority string is not 2 characters 'Pn', provided string: " + priority)
	}

	//ASCII value of '0' is decimal 48 and '9' is decimal 57.
	if priority[0] != 'P' || priority[1] < 49 || priority[1] > 53 {
		log.Fatal("priority myst be formatted as 'Pn' where P is capital and n is a number between 1 and 5, provided string: " + priority)
	}

	//Convert SHA to short-SHA (first 7 characters)
	var githubSha = os.Getenv("GITHUB_SHA")
	if len(githubSha) < 7 {
		githubSha = "(unknown)"
	} else {
		githubSha = githubSha[0:7]
	}

	//Create the client
	incidentClient, err := incident.NewClient(
		&client.Config{
			ApiKey: allArguments[2],
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	var falseType = false
	//var trueType = true

	//Create the incident
	_, err = incidentClient.Create(
		nil, &incident.CreateRequest{
			Message:     "Notification from: " + os.Getenv("GITHUB_REPOSITORY"),
			Description: "This incident is generated by the GitHub workflow: " + os.Getenv("GITHUB_WORKFLOW") + "\ncommit: " + githubSha + "\n\n" + allArguments[0],
			Responders: []incident.Responder{
				{Type: incident.Team, Name: allArguments[3]},
			},
			Tags: []string{allArguments[4]},
			Details: map[string]string{
				"branch":      os.Getenv("GITHUB_REF_NAME"),
				"runner arch": os.Getenv("RUNNER_ARCH"),
				"runner os":   os.Getenv("RUNNER_OS"),
			},
			Priority:  incident.Priority(priority),
			Note:      "Alert note",
			ServiceId: "github_action",
			StatusPageEntity: &incident.StatusPageEntity{
				Title:       "report from GitHub action",
				Description: "GitHub action reported an incident",
			},
			NotifyStakeholders: &falseType,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
