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

type Bid struct {
	Amount int    `json:"amount"`
	Item   string `json:"item,ommitempty"`
	BidID  string `json:"bid_id,ommitempty"`
	Email  string `json:"email,ommitempty"`
}

type UpdateBidInput struct {
	Amount int `json:"amount"`
}

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

// TODO: right now for all the bid stuff Im fetching the user based off their session. Think of a better way to do this.
func (r *Routes) PlaceBid(w http.ResponseWriter, req *http.Request) {
	log.Println("routes: attempting to place a new bid")

	// Get the item ID from the URL
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
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	if (item == repo.Item{}) {
		log.Printf("routes: could not find item %s \n", itemID)
		http.Error(w, "could not find item", 404)
		return
	}

	var in Bid
	if err := json.NewDecoder(req.Body).Decode(&in); err != nil {
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

	// check to see if this is now the top bid
	bid, err := r.Repo.GetBid(item.TopBid)
	if err != nil {
		log.Printf("routes: error getting top bid %v", err)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	// TODO: figure out a godo way to deal with the race condition that two different people can be updating to the top bid
	// update the top bid if this is the new top bid
	if in.Amount > bid.Amount {
		if err := r.Repo.UpdateItemsTopBid(in.BidID, itemID); err != nil {
			log.Printf("routes: error trying to update top bid %v \n", err)
			http.Error(w, internalErrorResponse, http.StatusInternalServerError)
			return
		}
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
		log.Println("routes: User %s was trying to update a bid they did not make", user.Email)
		http.Error(w, internalErrorResponse, http.StatusForbidden)
		return
	}

	if err := r.Repo.UpdateBid(id, in.Amount); err != nil {
		log.Println("routes: Error trying to update bid")
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	// get the item and the top bid to see if this is the new top bid
	item, err := r.Repo.GetItem(bid.ItemID)
	if err != nil {
		log.Printf("routes: Error trying to get item %s \n", bid.ItemID)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}
	topBid, err := r.Repo.GetBid(item.TopBid)
	if err != nil {
		log.Printf("routes: Error getting bid based off id %s \n", id)
		http.Error(w, internalErrorResponse, http.StatusInternalServerError)
		return
	}

	if in.Amount > topBid.Amount {
		if err := r.Repo.UpdateItemsTopBid(bid.BidID, item.ID); err != nil {
			log.Printf("routes: error trying to update top bid %v \n", err)
			http.Error(w, internalErrorResponse, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
