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
	TopBid      string `dynamodbav:"top_bid"`
}

func (r *Repo) UpdateItemsTopBid(topBid, itemID string) error {
	log.Printf("repo: Attempting to update the top bid for %s \n", itemID)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":tp": {
				S: aws.String(topBid),
			},
		},
		TableName: aws.String("millennium-falcon-auction-items"),
		Key: map[string]*dynamodb.AttributeValue{
			"itemID": {
				S: aws.String(itemID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set top_bid = :tp"),
	}

	_, err := r.svc.UpdateItem(input)
	if err != nil {
		log.Printf("repo: Error attempting to update item: %v", err)
		return err
	}

	log.Println("repo: Successfully updated item in repo")
	return nil
}

func (r *Repo) GetItem(itemID string) (Item, error) {
	log.Printf("repo: Getting item. \n", itemID)

	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("millennium-falcon-auction-items"),
		Key: map[string]*dynamodb.AttributeValue{
			"itemID": {
				S: aws.String(itemID),
			},
		},
	})
	if err != nil {
		return Item{}, errors.New("could not retrieve item from dynamo")
	}

	log.Println("Successfully retrieved item from dynamo.")

	item := Item{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &item); err != nil {
		return Item{}, err
	}
	return item, nil
}
