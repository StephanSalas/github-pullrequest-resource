package resource

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type GitHub interface {
	CreatePullRequest(newPullRequest *github.NewPullRequest) (pullRequest *github.PullRequest, responseCode int, err error)
}

type GitHubClient struct {
	client      *github.Client
	owner       string
	repository  string
	accessToken string
}

func NewGitHubClient(source Source) (*GitHubClient, error) {
	var httpClient = &http.Client{}
	var ctx = context.TODO()

	if source.Insecure {
		httpClient.Transport = &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	}

	if source.AccessToken != "" {
		var err error
		httpClient, err = oauthClient(ctx, source)
		if err != nil {
			return nil, err
		}
	}

	client := github.NewClient(httpClient)

	if source.GitHubAPIURL != "" {
		var err error
		if !strings.HasSuffix(source.GitHubAPIURL, "/") {
			source.GitHubAPIURL += "/"
		}
		client.BaseURL, err = url.Parse(source.GitHubAPIURL)
		if err != nil {
			return nil, err
		}

		client.UploadURL, err = url.Parse(source.GitHubAPIURL)
		if err != nil {
			return nil, err
		}
	}

	return &GitHubClient{
		client:      client,
		repository:  source.Repository,
		accessToken: source.AccessToken,
	}, nil
}

func oauthClient(ctx context.Context, source Source) (*http.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: source.AccessToken,
	})

	oauthClient := oauth2.NewClient(ctx, ts)

	githubHTTPClient := &http.Client{
		Transport: oauthClient.Transport,
	}

	return githubHTTPClient, nil
}

func (g *GitHubClient) CreatePullRequest(newPullRequest *github.NewPullRequest) (pullRequest *github.PullRequest, responseCode int, err error) {
	pullRequest, response, err := g.client.PullRequests.Create(context.TODO(), g.owner, g.repository, newPullRequest)
	// req.Header.Set("Accept", "application/octet-stream")
	// if g.accessToken != "" && req.URL.Host == g.client.BaseURL.Host {
	// 	req.Header.Set("Authorization", "Bearer "+g.accessToken)
	// }

	if err != nil {
		return pullRequest, response.StatusCode, err
	}

	return pullRequest, response.StatusCode, nil
}
