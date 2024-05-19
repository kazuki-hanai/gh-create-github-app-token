package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/config"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/github"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/jwt"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/token"

	"github.com/rs/zerolog"
)

func getLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return &logger
}

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

func main() {
	logger := getLogger()
	ctx := logger.WithContext(context.Background())

	config, err := config.NewConfig("", "", "")
	if err != nil {
		panic(err)
	}

	jwtToken, err := getJwtToken(config)
	if err != nil {
		panic(err)
	}

	token, err := getGitHubAppToken(ctx, config, jwtToken)
	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
