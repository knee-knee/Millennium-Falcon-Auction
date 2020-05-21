package repo

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
