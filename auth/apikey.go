package auth

import (
	"log"
	"os"
)

var (
	MatchServiceAPIKey = os.Getenv("MATCH_SERVICE_API_KEY")
	MatchServiceURL    = os.Getenv("MATCH_SERVICE_URL")
)

func init() {
	if MatchServiceURL == "" {
		log.Fatal("Environment variable MATCH_SERVICE_URL is undefined")
	}
	if MatchServiceAPIKey == "" {
		log.Fatal("Environment variable MATCH_SERVICE_API_KEY is undefined")
	}
}
