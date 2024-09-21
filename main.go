package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v65/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
)

func Hello(name string) string {
	return "Hello " + name
}

type Reviewer struct {
	Name  string
	Email string
}

type Reward struct {
	users     []string
	hash      string
	processed bool
}

var client *github.Client
var reviewer []Reviewer
var rewards []Reward

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
	return user.GetLogin(), nil
}

func getOpenPRs(ctx context.Context) {
	username, err := getAuthenticatedUsername(ctx)
	if err != nil {
		log.Fatalf("Failed to retrieve authenticated user's username: %v", err)
	}

	// Search for open pull requests authored by the authenticated user
	query := fmt.Sprintf("is:pr is:open author:%s", username)
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// Perform the search
	result, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		log.Fatalf("Error searching for pull requests: %v", err)
	}

	// Print the results
	log.Printf("Found %d open pull requests\n", result.GetTotal())
	for _, issue := range result.Issues {
		repoUrl := issue.GetRepositoryURL()
		repo := strings.Split(repoUrl, "/")
		repoNumber := issue.GetNumber()
		repoOwner := repo[len(repo)-2]
		repoName := repo[len(repo)-1]

		reviews := getAllReviews(ctx, repoOwner, repoName, repoNumber)
		if reviews == nil {
			return
		}

		// TODO: remove prints
		fmt.Printf("PR Title: %s\n", issue.GetTitle())
		fmt.Printf("Repo: %s\n", issue.GetHTMLURL())
		fmt.Printf("Repo: %v\n", issue.GetNumber())
		fmt.Printf("Owner: %s\n", repoOwner)
		fmt.Printf("Name: %s\n", repoName)
		for _, review := range reviews {
			user := review.GetUser().GetLogin()
			email := review.GetUser().GetEmail()
			fmt.Printf("reviewer: %s\n", user)
			fmt.Printf("email: %s\n", email)
		}
		fmt.Println("------")
	}
}

func getAllReviews(ctx context.Context, owner string, repo string, id int) []*github.PullRequestReview {
	opts := &github.ListOptions{PerPage: 50}

	reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, id, opts)
	if err != nil {
		log.Fatalln("Error getting PR reviewers:", err)
	}
	return reviews
}

func main() {
	initGitHubClient()
	ctx := context.Background()
	getOpenPRs(ctx)
}
