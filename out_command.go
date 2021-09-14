package resource

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/google/go-github/v39/github"
)

type OutCommand struct {
	github GitHub
	writer io.Writer
}

func NewOutCommand(github GitHub, writer io.Writer) *OutCommand {
	return &OutCommand{
		github: github,
		writer: writer,
	}
}

func (c *OutCommand) Run(request OutRequest) (outResponse OutResponse, err error) {
	var oResponse OutResponse
	gitHubClient, err := NewGitHubClient(request.Source)
	if err != nil {
		return oResponse, err
	}

	newPullRequest := &github.NewPullRequest{
		Title:               &request.Params.Title,
		Head:                &request.Params.SourceBranch,
		Base:                &request.Params.TargetBranch,
		Issue:               &request.Params.IssueRef,
		MaintainerCanModify: &request.Params.MaintainerCanModify,
		Draft:               &request.Params.isDraft,
	}

	createdPullRequest, responseCode, err := gitHubClient.CreatePullRequest(newPullRequest)

	if responseCode != 201 {
		return oResponse, fmt.Errorf("Received " + strconv.Itoa(responseCode) + "instead of expected 201")
	}

	oResponse.PullRequestLink = *createdPullRequest.URL
	return oResponse, nil
}

func (c *OutCommand) fileContents(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
