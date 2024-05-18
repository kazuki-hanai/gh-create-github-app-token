package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v62/github"
)

type GitHubAppTokenGenerator struct {
	PrivateKey   string
	AppID        string
	GitHubClient *github.Client
}

type GitHubAppTokenGenerateOptions struct {
	Org string
}

func (g *GitHubAppTokenGenerator) GenerateToken(options GitHubAppTokenGenerateOptions) (string, error) {
	jwtTokenGenerator := JwtTokenGenerator{
		PrivateKey: g.PrivateKey,
	}

	now := time.Now().Unix()
	// Set issue time to 30 seconds in the past to account for clock drift
	nowWithSafetyMargin := now - 30
	// JWT expiration time (10 minutes maximum)
	expiration := nowWithSafetyMargin + 60*10

	payload := JwtPayload{
		Exp: expiration,
		Iss: g.AppID,
		Iat: nowWithSafetyMargin,
	}

	jwtToken, err := jwtTokenGenerator.GenerateToken(payload)
	if err != nil {
		return "", err
	}

	// Get the installation ID
	var v map[string]interface{}
	req, err := g.GitHubClient.NewRequest("GET", fmt.Sprintf("/orgs/%s/installation", options.Org), nil)
	if err != nil {
		return "", err
	}
	_, err = g.GitHubClient.WithAuthToken(jwtToken).Do(context.Background(), req, &v)
	if err != nil {
		return "", err
	}

	installationID := int64(v["id"].(float64))

	installationToken, _, err := g.GitHubClient.WithAuthToken(jwtToken).Apps.CreateInstallationToken(context.Background(), installationID, nil)
	if err != nil {
		return "", err
	}

	return installationToken.GetToken(), nil
}
