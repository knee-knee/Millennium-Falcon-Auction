package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/millennium-falcon-auction/repo"
)

const authHeader = "auth"

// Bid represents a bid for an item.
type Bid struct {
	Amount int    `json:"amount"`
	Item   string `json:"item,ommitempty"`
	BidID  string `json:"bid_id,ommitempty"`
	Email  string `json:"email,ommitempty"`
}

// UpdateBidInput is the input to update a bid.
type UpdateBidInput struct {
	Amount int `json:"amount"`
}

// PlaceBidOutput is the output when you place a bid.
type PlaceBidOutput struct {
	BidID string `json:"bid_id"`
}

func (b Bid) toDyanmo() repo.Bid {
	return repo.Bid{
		Amount: b.Amount,
		Bidder: b.Email,
		ItemID: b.Item,
		BidID:  b.BidID,
	}
}

// PlaceBid will place a bid on an item.
// TODO: right now for all the bid stuff Im fetching the user based off their session. Think of a better way to do this.
func (r *Routes) PlaceBid(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: attempting to place a new bid")

	// Get the item ID from the URL
	params := mux.Vars(req)
	itemID, ok := params["item_id"]
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
	log.Printf("routes: Succesfully got item info for %s \n", itemID)

	var in Bid
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		log.Printf("routes: Error trying to decode body %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
	in.BidID = uuid.New().String()
	in.Item = itemID

	session := req.Header.Get(authHeader)
	user, err := r.Repo.GetUserBySession(session)
	if err != nil {
		log.Println("routes: Errror getting the user based of their session")
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	in.Email = user.Email
	if err := r.Repo.CreateBid(in.toDyanmo()); err != nil {
		log.Printf("routes: error creating bid %v", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	// Check to see if this is now the highest bid.
	if err := r.UpdateHighestBid(itemID); err != nil {
		log.Printf("routes: Error trying to check and see if this is the highest bid %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	out := PlaceBidOutput{
		BidID: in.BidID,
	}
	body, err := json.Marshal(out)
	if err != nil {
		log.Printf("routes: Error marshalling in to response body %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

// UpdateBid will allow a user to update the amount on their bid.
func (r *Routes) UpdateBid(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: attempting to update an existing bid")

	params := mux.Vars(req)
	id, ok := params["id"]
	if !ok {
		log.Println("routes: Request made without bid ID")
		http.Error(w, "did not provide bid ID", 400)
		return
	}

	var in UpdateBidInput
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
		log.Printf("routes: error decoding input %v", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	bid, err := r.Repo.GetBid(id)
	if err != nil {
		log.Printf("routes: Error getting bid based off id %s \n", id)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	session := req.Header.Get(authHeader)
	user, err := r.Repo.GetUserBySession(session)
	if err != nil {
		log.Println("routes: Errror getting the user based of their session")
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	if bid.Bidder != user.Email {
		log.Printf("routes: User %s was trying to update a bid they did not make \n", user.Email)
		http.Error(w, "Cannot update a Bid You Did Not Make", http.StatusForbidden)
		return
	}

	updatedBid, err := r.Repo.UpdateBid(id, in.Amount)
	if err != nil {
		log.Println("routes: Error trying to update bid")
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	if err := r.UpdateHighestBid(updatedBid.ItemID); err != nil {
		log.Printf("routes: Error trying to check and see if this is the highest bid %v \n", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	log.Printf("routes: Succesfully updated bid %s \n", id)

	w.WriteHeader(http.StatusNoContent)
}
