package token

import (
	"context"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/config"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/github"
)

type GitHubAppTokenGenerator interface {
	GenerateToken(ctx context.Context) (string, error)
}

type GitHubAppTokenGeneratorImpl struct {
	PrivateKey   string
	AppID        string
	Org          string
	GitHubClient *github.GitHubClient
}

func NewGitHubAppTokenGeneratorImpl(c *config.Config, gitHubClient *github.GitHubClient) (GitHubAppTokenGenerator, error) {
	return &GitHubAppTokenGeneratorImpl{
		PrivateKey:   c.PrivateKey,
		AppID:        c.AppID,
		Org:          c.Org,
		GitHubClient: gitHubClient,
	}, nil
}

func (g *GitHubAppTokenGeneratorImpl) GenerateToken(ctx context.Context) (string, error) {
	installationID, err := (*g.GitHubClient).GetInstallationID(ctx, g.Org)
	if err != nil {
		return "", err
	}

	installationToken, err := (*g.GitHubClient).CreateInstallationToken(ctx, installationID)
	if err != nil {
		return "", err
	}

	return installationToken, nil
}
