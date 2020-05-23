package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/millennium-falcon-auction/repo"
)

const numberOfBidsToDisplay = 10

type ItemInfo struct {
	Description string `json:"description"`
	Bids        []Bid  `json:"Bids"`
}

func bidsFromRepo(b []repo.Bid) []Bid {
	bids := make([]Bid, len(b))
	for i, bid := range b {
		bids[i] = Bid{
			Amount: bid.Amount,
			Email:  bid.Bidder,
			Item:   bid.ItemID,
			BidID:  bid.BidID,
		}
	}

	return bids
}

func (r *Routes) GetItemInfo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	itemID, ok := params["id"]
	if !ok {
		log.Println("routes: Request made without itemID")
		http.Error(w, "did not provide item ID", http.StatusBadRequest)
		return
	}
	log.Printf("routes: Getting item info for %s \n", itemID)

	item, err := r.Repo.GetItem(itemID)
	if err != nil {
		log.Printf("routes: Error getting item: %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	if (item == repo.Item{}) {
		log.Printf("routes: could not find item %s \n", itemID)
		http.Error(w, "could not find item", http.StatusNotFound)
		return
	}

	log.Printf("routes: Attempting to get %d bids for %s", numberOfBidsToDisplay, itemID)
	bids, err := r.Repo.GetTopBids(itemID, numberOfBidsToDisplay)
	if err != nil {
		log.Printf("routes: Error getting bids for %s: %v \n", itemID, err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	itemInfo := ItemInfo{
		Description: item.Description,
		Bids:        bidsFromRepo(bids),
	}

	body, err := json.Marshal(itemInfo)
	if err != nil {
		log.Println("routes: could not marshal into response")
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	w.Write(body)
}
