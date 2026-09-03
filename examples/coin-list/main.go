// coin-list retrieves a short list of coins sorted by Galaxy Score.
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

	sort, limit, descending := "galaxy_score", 10, true
	client := lunarcrush.NewClient(apiKey, lunarcrush.WithTimeout(15*time.Second))
	coins, err := client.Coins.List(context.Background(), &lunarcrush.CoinsListParams{
		Sort:  &sort,
		Limit: &limit,
		Desc:  &descending,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, coin := range coins.Data {
		fmt.Printf("%-8s score=%5.1f price=$%.4f\n", coin.Symbol, coin.GalaxyScore, coin.Price)
	}
}
