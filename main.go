package main

import (
	"fmt"
	"os"

	"github.com/kazuki-hanai/gh-create-github-app-token/cmd"
	"github.com/kazuki-hanai/gh-create-github-app-token/pkg/logger"
)

func main() {
	logger.ConfigureLogger()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
