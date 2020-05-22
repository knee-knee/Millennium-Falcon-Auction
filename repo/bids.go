package repo

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Bid struct {
	Amount int    `dynamodbav:"amount"`
	Bidder string `dynamodbav:"bidder_email"`
	ItemID string `dynamodbav:"item_id"`
	BidID  string `dynamodbav:"bidID"`
}

func (r *Repo) CreateBid(in Bid) error {
	log.Println("repo: attempting to create a new bid in dyanmo.")
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		return errors.New("repo: could not marshal created question into dynamo map")
	}

	if _, err = r.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("millennium-falcon-auction-bids"),
		Item:      item,
	}); err != nil {
		return errors.New("repo: could not put created bid into dynamo")
	}

	log.Printf("repo: successfully created bid %s", in.BidID)

	return nil
}

func (r *Repo) GetBid(bidID string) (Bid, error) {
	log.Printf("repo: Getting bid %s", bidID)

	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("millennium-falcon-auction-bids"),
		Key: map[string]*dynamodb.AttributeValue{
			"bidID": {
				S: aws.String(bidID),
			},
		},
	})
	if err != nil {
		log.Printf("repo: Error getting bid from dyanmo %v \n", err)
		return Bid{}, errors.New("could not retrieve bid from dynamo")
	}

	log.Println("repo: Successfully retrieved bid from dyanmo.")

	var bid Bid
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &bid); err != nil {
		return Bid{}, err
	}
	return bid, nil
}

func (r *Repo) GetTopBids(itemID string, numberOfBids int) ([]Bid, error) {
	log.Printf("repo: Getting %d bids for %s", numberOfBids, itemID)

	resp, err := r.svc.Query(&dynamodb.QueryInput{
		TableName: aws.String("millennium-falcon-auction-bids"),
		IndexName: aws.String("item_id-amount-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"item_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(itemID),
					},
				},
			},
		},
		ScanIndexForward: aws.Bool(true),
		Limit:            aws.Int64(int64(numberOfBids)),
	})
	if err != nil {
		return nil, err
	}
	if resp.Count == nil {
		log.Printf("repo: count from scan is empty \n")
		return nil, errors.New("count of scan was empty")
	}

	bids := make([]Bid, *resp.Count)
	if err := dynamodbattribute.UnmarshalListOfMaps(resp.Items, &bids); err != nil {
		log.Printf("repo: error unmarshaling map: %v \n", err)
		return nil, err
	}

	for i, bid := range bids {
		fmt.Printf("i: %d: bid amount: %d bid email: %s  \n", i, bid.Amount, bid.Bidder)
	}

	return bids, nil
}
