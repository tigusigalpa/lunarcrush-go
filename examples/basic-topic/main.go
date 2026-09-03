// basic-topic retrieves a topic's social summary.
package main

import (
	"context"
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

	client := lunarcrush.NewClient(apiKey, lunarcrush.WithTimeout(15*time.Second))
	topic, err := client.Topics.Get(context.Background(), "bitcoin")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s ranks #%d with %.0f interactions in 24 hours.\n",
		topic.Data.Title, topic.Data.TopicRank, topic.Data.Interactions24h)
}
