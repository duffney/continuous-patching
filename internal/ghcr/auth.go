package ghcr

import (
	"log"

	"github.com/docker/cli/cli/config"
)

func GetAuthToken() string {
	// Load Docker config
	cfg, err := config.Load(config.Dir())
	if err != nil {
		log.Fatalf("Error loading Docker config: %v", err)
	}

	// Get authentication for ghcr.io
	authConfig, err := cfg.GetAuthConfig("ghcr.io")
	if err != nil {
		log.Fatalf("Error getting auth config: %v", err)
	}

	// Determine which token to use
	var token string
	if authConfig.IdentityToken != "" {
		token = authConfig.IdentityToken
	} else if authConfig.Password != "" {
		token = authConfig.Password // For GitHub, the password is often a PAT
	} else {
		//TODO: add support for env token
		log.Fatal("No valid token found for ghcr.io in Docker config")
	}
	return token
}
