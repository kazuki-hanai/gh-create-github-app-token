package main

import (
	"fmt"
	"os"

	"github.com/google/go-github/v62/github"
	"github.com/kazuki-hanai/create-github-app-token/internal"
)


func main() {
	org := os.Getenv("ORG")
	privateKey := os.Getenv("PRIVATE_KEY")
	appID := os.Getenv("APP_ID")

	g := internal.GitHubAppTokenGenerator{
		PrivateKey: privateKey,
		AppID: appID,
		GitHubClient: github.NewClient(nil),
	}

	token, err := g.GenerateToken(internal.GitHubAppTokenGenerateOptions{
		Org: org,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
