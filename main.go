package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v65/github"
	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Hello(name string) string {
	return "Hello " + name
}

type Reviewer struct {
	Username string
	Name     string
	Email    string
}

type Reward struct {
	users     []Reviewer
	hash      string
	processed bool
}

var client *github.Client
var rewards []Reward
var ctx = context.Background()

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

// generates a hash
func hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
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
		reviewers := make([]Reviewer, 0)

		reviews := getAllReviews(ctx, repoOwner, repoName, repoNumber)
		if reviews == nil {
			return
		}

		hashInput := fmt.Sprintf("%v-%v-%v", repoOwner, repoName, repoNumber)
		prHash := fmt.Sprintf("%v", hash(hashInput))

		for _, review := range reviews {
			user := review.GetUser().GetLogin()
			// do not include self in the reviewers list
			if user == username {
				continue
			}
			name := review.GetUser().GetName()
			//email := review.GetUser().GetEmail()
			reviewers = append(reviewers, Reviewer{
				Username: user,
				Name:     name,
				Email:    "farhan.alam@optimizely.com",
				//Email: email,
			})
		}
		rewards = append(rewards, Reward{
			users:     reviewers,
			hash:      prHash,
			processed: false,
		})
	}
}

// TODO: force get email
func forceGetEmail(ctx context.Context, user string, owner string, repo string) string {
	return "faa@faa.com"
}

func getAllReviews(ctx context.Context, owner string, repo string, id int) []*github.PullRequestReview {
	opts := &github.ListOptions{PerPage: 50}

	reviews, _, err := client.PullRequests.ListReviews(ctx, owner, repo, id, opts)
	if err != nil {
		log.Fatalln("Error getting PR reviewers:", err)
	}
	return reviews
}
func removeDuplicateUsers(users []Reviewer) []Reviewer {
	temp := make(map[string]bool)
	var result []Reviewer

	for _, user := range users {
		if !temp[user.Username] {
			temp[user.Username] = true
			result = append(result, user)
		}
	}
	return result
}

func processRewardList() {
	reviewers := make([]Reviewer, 0)
	for _, value := range rewards {
		reviewers = append(reviewers, value.users...)
	}
	reviewers = removeDuplicateUsers(reviewers)
	log.Printf("Found %v unique reviewers\n", len(reviewers))
	// Force get emails here

	usernames := make([]string, 0)
	for _, reviewer := range reviewers {
		if reviewer.Email == "" {
			// use Bonusly autocomplete as a last resort
		}
		username, err := getBonuslyUsernames(reviewer.Email)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		usernames = append(usernames, username)
		fmt.Printf("username %v\n", username)
	}
	message := generateBonuslyMessage(usernames)
	fmt.Printf("message %v\n", message)
}

func generateBonuslyMessage(usernames []string) string {
	tag := "focus-on-continuous-improvement"
	points := 5
	for i, name := range usernames {
		usernames[i] = "@" + name
	}
	userStr := strings.Join(usernames, " ")

	message := fmt.Sprintf("%s Thanks for reviewing my code +%d #%s", userStr, points, tag)

	return message
}

func getBonuslyUsernames(email string) (string, error) {
	token := os.Getenv("BONUSLY_ACCESS_TOKEN")
	encodedEmail := url.QueryEscape(email)
	requestUrl := fmt.Sprintf("https://bonus.ly/api/v1/users?limit=1&email=%v&include_archived=false", encodedEmail)
	req, _ := http.NewRequest("GET", requestUrl, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	name := gjson.Get(string(body), "result.0.username")
	if name.Exists() {
		return name.String(), nil
	} else {
		return "", fmt.Errorf("username not found")
	}
}

func sendBonuslyPoints(message string) {

}

func main() {
	initGitHubClient()
	getOpenPRs(ctx)
	processRewardList()
}
