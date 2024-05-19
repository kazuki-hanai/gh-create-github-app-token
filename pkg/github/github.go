package github

import (
	"context"

	"github.com/google/go-github/v62/github"
)

type GitHubClient interface {
	GetInstallationID(ctx context.Context, org string) (int64, error)
	CreateInstallationToken(ctx context.Context, installationID int64) (string, error)
}

type GitHubClientImpl struct {
	Client   *github.Client
	jwtToken string
}

func NewGitHubClientImpl(jwtToken string) (GitHubClient, error) {
	c := github.NewClient(nil)

	return &GitHubClientImpl{
		Client:   c,
		jwtToken: jwtToken,
	}, nil
}

func (g *GitHubClientImpl) GetInstallationID(ctx context.Context, org string) (int64, error) {
	req, err := g.Client.NewRequest("GET", "/orgs/"+org+"/installation", nil)
	if err != nil {
		return 0, err
	}

	var v map[string]interface{}
	_, err = g.Client.WithAuthToken(g.jwtToken).Do(ctx, req, &v)
	if err != nil {
		return 0, err
	}

	installationID := int64(v["id"].(float64))
	return installationID, nil
}

func (g *GitHubClientImpl) CreateInstallationToken(ctx context.Context, installationID int64) (string, error) {
	installationToken, _, err := g.Client.WithAuthToken(g.jwtToken).Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return "", err
	}
	return installationToken.GetToken(), nil
}
