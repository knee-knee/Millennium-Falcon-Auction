package routes

import (
	"log"
)

func (r *Routes) CheckAndUpdateIfBidIsHightest(amount int, bidID, itemID string) error {
	// Locking here because if two bids enter here higher than the current highest bid we need to lock
	r.highestBidMux.Lock()
	defer r.highestBidMux.Unlock()

	// get the item and the top bid to see if this is the new top bid
	item, err := r.Repo.GetItem(itemID)
	if err != nil {
		log.Printf("routes: Error trying to get item %s \n", itemID)
		return err
	}
	topBid, err := r.Repo.GetBid(item.TopBid)
	if err != nil {
		log.Printf("routes: Error getting bid based off id %s \n", item.TopBid)
		return err
	}

	if amount > topBid.Amount {
		if err := r.Repo.UpdateItemsTopBid(bidID, itemID); err != nil {
			log.Printf("routes: error trying to update top bid %v \n", err)
			return err
		}
	}

	return nil
}
