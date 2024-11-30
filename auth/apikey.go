package auth

import (
	"log"
	"os"
)

var (
	MatchServiceClientAPIKey = os.Getenv("MATCH_SERVICE_CLIENT_API_KEY")
	MatchServiceURL          = os.Getenv("MATCH_SERVICE_URL")
)

func init() {
	if MatchServiceURL == "" {
		log.Fatal("Environment variable MATCH_SERVICE_URL is undefined")
	}
	if MatchServiceClientAPIKey == "" {
		log.Fatal("Environment variable MATCH_SERVICE_CLIENT_API_KEY is undefined")
	}
}
