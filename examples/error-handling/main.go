// error-handling shows retry configuration and API error handling.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	lunarcrush "github.com/tigusigalpa/lunarcrush-go"
)

func main() {
	apiKey := os.Getenv("LUNARCRUSH_API_KEY")
	if apiKey == "" {
		log.Fatal("set LUNARCRUSH_API_KEY before running this example")
	}

	client := lunarcrush.NewClient(apiKey,
		lunarcrush.WithTimeout(15*time.Second),
		lunarcrush.WithRetry(3, time.Second),
	)
	_, err := client.Topics.Get(context.Background(), "bitcoin")
	if err == nil {
		fmt.Println("request succeeded")
		return
	}

	switch {
	case errors.Is(err, lunarcrush.ErrUnauthorized):
		log.Fatal("the API key is invalid or missing")
	case errors.Is(err, lunarcrush.ErrNotFound):
		log.Fatal("the requested topic does not exist")
	case errors.Is(err, lunarcrush.ErrRateLimited):
		log.Fatal("the request remained rate limited after retries")
	default:
		log.Fatalf("request failed: %v", err)
	}
}
