package repo

import (
	"fmt"
	"log"
)

type Bid struct {
	Amount int
	Bidder string
}

func (r *Repo) GetItemDescription(itemID string) (string, error) {
	log.Printf("repo: Getting item description. \n", itemID)

	return "Some dumb item description placeholder", nil
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
