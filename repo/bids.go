package repo

import (
	"errors"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Bid is the DB represntation of a bid.
type Bid struct {
	Amount int    `dynamodbav:"amount"`
	Bidder string `dynamodbav:"bidder_email"`
	ItemID string `dynamodbav:"item_id"`
	BidID  string `dynamodbav:"bidID"`
}

const bidsTableID = "millennium-falcon-auction-bids"

// CreateBid will create a new bid in dyanmoDB.
func (r *Repo) CreateBid(in Bid) error {
	log.Println("repo: attempting to create a new bid in dyanmo.")
	item, err := dynamodbattribute.MarshalMap(in)
	if err != nil {
		log.Printf("repo: Error marshaling map %v \n", err)
		return err
	}

	if _, err = r.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(bidsTableID),
		Item:      item,
	}); err != nil {
		log.Printf("repo: Error putting in bid %v \n", err)
		return err
	}

	log.Printf("repo: successfully created bid %s \n", in.BidID)

	return nil
}

// DeleteBid will delete a bid in dyanmoDB.
func (r *Repo) DeleteBid(bid Bid) error {
	log.Printf("repo: attempting to delete bid %s in dyanmo \n", bid.BidID)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"bidID": {
				S: aws.String(bid.BidID),
			},
			"amount": {
				N: aws.String(strconv.Itoa(bid.Amount)),
			},
		},
		TableName: aws.String(bidsTableID),
	}

	_, err := r.svc.DeleteItem(input)
	if err != nil {
		log.Printf("repo: error trying to delete bid %v \n", err)
		return err
	}

	log.Printf("repo: Succesfully deleted bid %s \n", bid.BidID)

	return nil
}

// UpdateBid is used to update an existing bid.
// This will only allow you to update the amount on a bid.
// The steps involved in updating a a bid are: first get the bid, delete it, then create a new entry with the updated amount.
func (r *Repo) UpdateBid(bidID string, amount int) (Bid, error) {
	log.Printf("repo: Attempting to update bid %v with new amount %d", bidID, amount)

	// first get the bid
	bid, err := r.GetBid(bidID)
	if err != nil {
		log.Println("repo: error getting bid")
		return Bid{}, err
	}

	// delete bid
	if err := r.DeleteBid(bid); err != nil {
		log.Println("repo: error deleting bid")
		return Bid{}, err
	}

	// change bid to updated values
	bid.Amount = amount

	// create new bid
	if err := r.CreateBid(bid); err != nil {
		log.Println("repo: error creating bid")
		return Bid{}, err
	}

	log.Printf("repo: Succesfully created bid %s \n", bidID)
	return bid, nil
}

// GetBid will get a big based off the bid ID.
func (r *Repo) GetBid(bidID string) (Bid, error) {
	log.Printf("repo: Getting bid %s", bidID)

	resp, err := r.svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(bidsTableID),
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

// GetTopBids will return a certain amount of top bids for a certain item.
func (r *Repo) GetTopBids(itemID string, numberOfBids int) ([]Bid, error) {
	log.Printf("repo: Getting %d bids for %s", numberOfBids, itemID)

	resp, err := r.svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(bidsTableID),
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
		log.Println("repo: Error querying for top bids ")
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

	log.Printf("repo: Succesfully retrived top bids for item %s \n", itemID)
	return bids, nil
}
