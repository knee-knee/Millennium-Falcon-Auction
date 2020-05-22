package repo

import (
	"errors"
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

// TODO: figure out why get item is not working here
func (r *Repo) GetBid(bidID string) (Bid, error) {
	log.Printf("repo: Getting bid %s", bidID)

	resp, err := r.svc.Query(&dynamodb.QueryInput{
		TableName: aws.String("millennium-falcon-auction-bids"),
		IndexName: aws.String("bidID-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"bidID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(bidID),
					},
				},
			},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int64(1),
	})
	if err != nil {
		return Bid{}, err
	}
	if resp.Count == nil {
		log.Printf("repo: count from scan is empty \n")
		return Bid{}, errors.New("count of scan was empty")
	}

	log.Println("repo: Successfully retrieved bid from dyanmo.")

	var bid Bid
	// TODO: this isn't great but its what is needed because we need to use the query.
	if err := dynamodbattribute.UnmarshalMap(resp.Items[0], &bid); err != nil {
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
		ScanIndexForward: aws.Bool(false),
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

	return bids, nil
}
