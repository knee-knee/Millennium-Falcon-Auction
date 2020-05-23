package routes

import (
	"log"
)

// UpdateHighestBid will take in the itemID and update the highest bid field to the current highest bid.
func (r *Routes) UpdateHighestBid(itemID string) error {
	// Locking here because if two bids enter here higher than the current highest bid we need to lock
	r.highestBidMux.Lock()
	defer r.highestBidMux.Unlock()

	topBids, err := r.Repo.GetTopBids(itemID, 1)
	if err != nil {
		log.Printf("repo: Error getting top two bids %v \n", err)
		return err
	}

	// is unnescary in the case where the top bid stays the same but for simplicity ill leave it like this.
	if err := r.Repo.UpdateItemsTopBid(topBids[0].BidID, itemID); err != nil {
		log.Printf("routes: error trying to update top bid %v \n", err)
		return err
	}

	return nil
}
