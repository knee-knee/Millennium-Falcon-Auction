package repo

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Item is the DB representaion of an item up for auction.
type Item struct {
	ID          string `dynamodbav:"ID"`
	Description string `dynamodbav:"description"`
	TopBid      string `dynamodbav:"top_bid"`
}

const itemsTableID = "millennium-falcon-auction-items"

// UpdateItemsTopBid will update the top bid ID for an item.
func (r *Repo) UpdateItemsTopBid(topBid, itemID string) error {
	log.Printf("repo: Attempting to update the top bid for %s \n", itemID)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":tp": {
				S: aws.String(topBid),
			},
		},
		TableName: aws.String(itemsTableID),
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

// GetItem will get an item based off the item ID.
func (r *Repo) GetItem(itemID string) (Item, error) {
	log.Printf("repo: Getting item. \n", itemID)

	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(itemsTableID),
		Key: map[string]*dynamodb.AttributeValue{
			"itemID": {
				S: aws.String(itemID),
			},
		},
	})
	if err != nil {
		log.Printf("repo: Error getting item from dyanmo %v \n", err)
		return Item{}, err
	}

	log.Println("Successfully retrieved item from dynamo.")

	item := Item{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &item); err != nil {
		log.Printf("repo: Error unmarshaling output into intem %v \n", err)
		return Item{}, err
	}

	log.Printf("repo: Succesfully retrieved item %s \n", itemID)
	return item, nil
}
