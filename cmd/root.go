package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/config"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/github"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/jwt"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/token"
)

var cfg config.Config

func getJwtToken(c *config.Config) (string, error) {
	jwtTokenGenerator, err := jwt.NewJwtTokenGenerator(c.PrivateKey)
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	// Set issue time to 30 seconds in the past to account for clock drift
	nowWithSafetyMargin := now - 30
	// JWT expiration time (10 minutes maximum)
	expiration := nowWithSafetyMargin + 60*10

	payload := jwt.JwtPayload{
		Exp: expiration,
		Iss: c.AppID,
		Iat: nowWithSafetyMargin,
	}

	jwtToken, err := jwtTokenGenerator.GenerateJwtToken(payload)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func getGitHubAppToken(ctx context.Context, c *config.Config, jwtToken string) (string, error) {
	gitHubClient, err := github.NewGitHubClientImpl(jwtToken)
	if err != nil {
		return "", err
	}

	gitHubAppTokenGeneratorImpl, err := token.NewGitHubAppTokenGeneratorImpl(c, &gitHubClient)
	if err != nil {
		return "", err
	}

	token, err := gitHubAppTokenGeneratorImpl.GenerateToken(ctx)
	if err != nil {
		return "", err
	}

	return token, nil
}

var RootCmd = &cobra.Command{
	Use:   "gh-create-github-app-token",
	Short: "gh-create-github-app-token is a extension of gh command to create a GitHub App token.",
	Long:  `gh-create-github-app-token is a extension of gh command to create a GitHub App token. It uses the GitHub App's private key to generate a JWT token and then uses the JWT token to generate a GitHub App token.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		jwtToken, err := getJwtToken(&cfg)
		if err != nil {
			panic(err)
		}

		token, err := getGitHubAppToken(ctx, &cfg, jwtToken)
		if err != nil {
			panic(err)
		}

		fmt.Println(token)
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&cfg.PrivateKey, "private-key", "p", "", "Path to the private key file")
	RootCmd.PersistentFlags().StringVarP(&cfg.AppID, "app-id", "a", "", "GitHub App ID")
	RootCmd.PersistentFlags().StringVarP(&cfg.Org, "org", "o", "", "GitHub App Installation ID")

	RootCmd.MarkPersistentFlagRequired("private-key")
	RootCmd.MarkPersistentFlagRequired("app-id")
}
