package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v65/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func Hello(name string) string {
	return "Hello " + name
}

var client *github.Client

// Initializes the GitHub client
func initGitHubClient() {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if token == "" {
		log.Fatalf("GITHUB_ACCESS_TOKEN not set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Initialize the global client
	client = github.NewClient(tc)
}

// Get the authenticated user's username
func getAuthenticatedUsername(ctx context.Context) (string, error) {
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return "", fmt.Errorf("error fetching authenticated user: %v", err)
	}
	fmt.Printf("Authenticated user perms: %s\n", user.Permissions)
	return user.GetLogin(), nil
}

func getAllMyOpenPullRequests(ctx context.Context) {
	// Get the authenticated username
	username, err := getAuthenticatedUsername(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve authenticated user's username: %v", err)
	}

	fmt.Printf("User name: %s\n", username)

	// Search for open pull requests authored by the authenticated user
	query := fmt.Sprintf("is:pr is:open author:%s", username)
	searchOpts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 10}, // Pagination: fetch 10 results per page
	}

	fmt.Printf("Search query: %s\n", query)

	// Perform the search
	result, _, err := client.Search.Issues(ctx, query, searchOpts)
	if err != nil {
		log.Fatalf("Error searching for pull requests: %v", err)
	}

	// Print the results
	fmt.Printf("Found %d open pull requests\n", result.GetTotal())
	for _, issue := range result.Issues {
		fmt.Printf("PR Title: %s\n", issue.GetTitle())
		fmt.Printf("PR URL: %s\n", issue.GetHTMLURL())
		fmt.Printf("Repo: %s\n", issue.GetRepositoryURL())
		fmt.Println("------")
	}
}

func main() {
	initGitHubClient()
	ctx := context.Background()
	getAllMyOpenPullRequests(ctx)
}
