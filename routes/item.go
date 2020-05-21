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
		bids[i].Amount = bid.Amount
		bids[i].Bidder = bid.Bidder
	}

	return bids
}

type Bid struct {
	Amount int    `json:"amount"`
	Bidder string `json:"bidder"`
}

func (r *Routes) GetItemInfo(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	itemID, ok := params["id"]
	if !ok {
		log.Println("routes: Request made without itemID")
		http.Error(w, "did not provide question ID", 400)
		return
	}
	log.Printf("routes: Getting item info for %s \n", itemID)

	des, err := r.Repo.GetItemDescription(itemID)
	if err != nil {
		log.Printf("routes: Error getting item description: %s \n", err)
		http.Error(w, "Internal Server Error", 500)
	}
	if des == "" {
		log.Printf("routes: user trying to search for item: %s which does not exists \n", itemID)
		http.Error(w, "could not find item", 404)
		return
	}

	log.Printf("routes: Attempting to get %d bids for %s", numberOfBidsToDisplay, itemID)
	bids, err := r.Repo.GetTopBids(itemID, numberOfBidsToDisplay)
	if err != nil {
		log.Printf("routes: Error getting bids for %s: %v \n", itemID, err)
		http.Error(w, "Internal Server Error", 500)
	}

	itemInfo := ItemInfo{
		Description: des,
		Bids:        bidsFromRepo(bids),
	}

	body, err := json.Marshal(itemInfo)
	if err != nil {
		log.Println("routes: could not marshal into response")
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Write(body)
}
