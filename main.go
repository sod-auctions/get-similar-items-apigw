package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sod-auctions/auctions-db"
	"log"
	"net/http"
	"os"
)

type ErrorMessage struct {
	Error string `json:"error"`
}

var database *auctions_db.Database

func init() {
	log.SetFlags(0)
	var err error
	database, err = auctions_db.NewDatabase(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
}

type Item struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	MediaURL string `json:"mediaUrl"`
	Rarity   string `json:"rarity"`
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := event.QueryStringParameters["name"]

	items, err := database.GetSimilarItems(name, 15)
	if err != nil {
		log.Printf("An error occurred: %v\n", err)

		errorMessage := ErrorMessage{Error: "An internal error occurred"}
		body, _ := json.Marshal(errorMessage)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "http://localhost:3000",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Authorization",
			},
			Body: string(body),
		}, nil
	}

	var mItems []*Item
	for _, item := range items {
		mItems = append(mItems, &Item{
			Id:       item.Id,
			Name:     item.Name,
			MediaURL: item.MediaURL,
			Rarity:   item.Rarity,
		})
	}

	body, _ := json.Marshal(mItems)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "http://localhost:3000",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept, Authorization",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
