package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/millennium-falcon-auction/repo"
)

type Bid struct {
	Email  string `json:"email"`
	Amount int    `json:"amount"`
	Item   string `json:"item,ommitempty"`
	BidID  string `json:"bid_id,ommitempty"`
}

type UpdateBidInput struct {
	Amount int `json:"amount"`
}

func (b Bid) toDyanmo() repo.Bid {
	return repo.Bid{
		Amount: b.Amount,
		Bidder: b.Email,
		ItemID: b.Item,
		BidID:  b.BidID,
	}
}

func (r *Routes) PlaceBid(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: attempting to place a new bid")

	params := mux.Vars(req)
	itemID, ok := params["item_id"]
	if !ok {
		log.Println("routes: Request made without itemID")
		http.Error(w, "did not provide item ID", 400)
		return
	}

	log.Printf("routes: Getting item info for %s \n", itemID)
	item, err := r.Repo.GetItem(itemID)
	if err != nil {
		log.Printf("routes: Error getting item: %v \n", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	if (item == repo.Item{}) {
		log.Printf("routes: could not find item %s \n", itemID)
		http.Error(w, "could not find item", 404)
		return
	}

	var in Bid
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer req.Body.Close()
	in.BidID = uuid.New().String()
	in.Item = itemID

	if err := r.Repo.CreateBid(in.toDyanmo()); err != nil {
		log.Printf("routes: error creating bid %v", err)
		http.Error(w, "Internal Sever Error", http.StatusInternalServerError)
		return
	}

	// check to see if this is now the top bid
	bid, err := r.Repo.GetBid(item.TopBid)
	if err != nil {
		log.Printf("routes: error getting top bid %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// update the top bid if this is the new top bid
	if in.Amount > bid.Amount {
		if err := r.Repo.UpdateItemsTopBid(in.BidID, itemID); err != nil {
			log.Printf("routes: error trying to update top bid %v \n", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

}

// Ensure here that the person that is updating the bid is the same user.
// Also im just going to assume all you can update is the amount
func (r *Routes) UpdateBid(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: attempting to update an existing bid")

	params := mux.Vars(req)
	id, ok := params["id"]
	if !ok {
		log.Println("routes: Request made without bid ID")
		http.Error(w, "did not provide bid ID", 400)
		return
	}

	// here is where I would get the email, but I will want to add auth to do that
	// Also check to see if this is now the top bid for the item

	var in UpdateBidInput
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		log.Printf("routes: error decoding input %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer req.Body.Close()

	if err := r.Repo.UpdateBid(id, in.Amount); err != nil {
		log.Println("routes: Error trying to update bid")
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
