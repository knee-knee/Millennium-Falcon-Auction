package repo

import (
	"errors"
	"log"
	"strconv"

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

func (r *Repo) DeleteBid(bid Bid) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"bidID": {
				S: aws.String(bid.BidID),
			},
			"amount": {
				N: aws.String(strconv.Itoa(bid.Amount)),
			},
		},
		TableName: aws.String("millennium-falcon-auction-bids"),
	}

	_, err := r.svc.DeleteItem(input)
	if err != nil {
		log.Printf("repo: error trying to delete bid %v \n", err)
		return err
	}

	return nil
}

// going to assume here that the bid you are updating is just the amount
// also since the amount is part of the key im going to make this a little hacky
// Im going to first get the bid, delete it, then create a new entry with the updated amount
func (r *Repo) UpdateBid(bidID string, amount int) error {
	log.Printf("repo: Attempting to update bid %v with new amount %d", bidID, amount)

	// first get the bid
	bid, err := r.GetBid(bidID)
	if err != nil {
		log.Println("repo: error getting bid")
		return err
	}

	// delete bid
	if err := r.DeleteBid(bid); err != nil {
		log.Println("repo: error deleting bid")
		return err
	}

	// change bid to updated values
	bid.Amount = amount

	// create new bid
	if err := r.CreateBid(bid); err != nil {
		log.Println("repo: error creating bid")
		return err
	}

	return nil
}

// TODO: figure out why I have to use query
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
		log.Printf("repo: count from query is empty \n")
		return Bid{}, errors.New("count of query was empty")
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
