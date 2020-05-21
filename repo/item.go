package repo

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Bid struct {
	Amount int
	Bidder string
}

type Item struct {
	ID          string `dynamodbav:"ID"`
	Description string `dynamodbav:"description"`
}

func (r *Repo) GetItemDescription(itemID string) (string, error) {
	log.Printf("repo: Getting item description. \n", itemID)

	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("millennium-falcon-auction-items"),
		Key: map[string]*dynamodb.AttributeValue{
			"itemID": {
				S: aws.String(itemID),
			},
		},
	})
	if err != nil {
		return "", errors.New("could not retrieve item from dynamo")
	}

	log.Println("Successfully retrieved item from dynamo.")

	item := Item{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &item); err != nil {
		return "", err
	}
	return item.Description, nil
}

func (r *Repo) GetTopBids(itemID string, numberOfBids int) ([]Bid, error) {
	log.Printf("repo: Getting %d bids for %s", numberOfBids, itemID)

	bids := make([]Bid, numberOfBids)

	for i := 0; i < numberOfBids; i++ {
		bids[i].Amount = i + 1
		bids[i].Bidder = fmt.Sprintf("Test User %d", i)
	}

	return bids, nil
}
