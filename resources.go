package resource

type Source struct {
	Title               string `json:"title"`
	GitHubAPIURL        string `json:"github_api_url"`
	Owner               string `json:"owner"`
	Repository          string `json:"repository"`
	SourceBranch        string `json:"sourceBranch"`
	TargetBranch        string `json:"targetBranch"`
	AccessToken         string `json:"access_token"`
	IssueRef            int    `json:"issue_ref"`
	MaintainerCanModify bool   `json:"maintainer_can_modify"`
	isDraft             bool   `json:"is_draft"`
	Insecure            bool   `json:"insecure"`
}

func NewOutRequest() OutRequest {
	res := OutRequest{}
	return res
}

type OutRequest struct {
	Source Source `json:"source"`
}

type OutResponse struct {
	PullRequestLink string `json:"pull_request_link"`
}
