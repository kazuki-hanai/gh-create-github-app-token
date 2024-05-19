package config

import (
	"fmt"
	"os"
)

type Config struct {
    PrivateKey   string
    AppID        string
	Org 		 string
}

func NewConfig(org, privateKey, appID string) (*Config, error) {
	if org == "" {
		org = os.Getenv("ORG")
	}

	if org == "" {
		return nil, fmt.Errorf("ORG is not set")
	}

	if privateKey == "" {
		privateKey = os.Getenv("PRIVATE_KEY")
	}

	if privateKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY is not set")
	}

	if appID == "" {
		appID = os.Getenv("APP_ID")
	}

	if appID == "" {
		return nil, fmt.Errorf("APP_ID is not set")
	}

	c := &Config{
		PrivateKey: privateKey,
		AppID: appID,
		Org: org,
	}
	return c, nil
}
