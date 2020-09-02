package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	gms "github.com/blackhatbrigade/gomessagestore"
	"github.com/blackhatbrigade/gomessagestore/uuid"
	log "github.com/sirupsen/logrus"
)

// GetHandler here's where the incoming request will be received
func GetHandler(repo Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {} // get does a count* for candidates and returns
}

// PostHandler is for our post requests
func PostHandler(repo Repository, ms gms.MessageStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// marshal to json and write a command to the messagestore
		// read in the request
		// example request []byte{`
		//  "candidate": "bob",
		//  "userId": "1",
		//`}
		var request map[string]interface{}
		json.NewDecoder(r.Body).Decode(&request)
		ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		// metadata
		metadata := []byte(fmt.Sprintf(
			`{"userId": "%s"}`,
			request["userId"]))

		// Return the response
		// set the header
		w.Header().Set("Content-Type", "application/json")

		// put response into a map
		response := map[string]interface{}{
			"candidate": request["candidate"],
			"userId":    request["userId"],
		}
		voteCommandData, err := json.Marshal(voteCmd{
			Candidate: request["candidate"].(string),
		})
		if err != nil {
			log.Error("no no no, you didn't json correctly", err)
		}
		userID := uuid.Must(uuid.Parse(request["userId"].(string)))

		// pack and write a command
		cmd := gms.Command{
			ID:             gms.NewID(),
			EntityID:       userID,
			StreamCategory: "Vote",
			MessageType:    "VoteForCandidate",
			Data:           voteCommandData,
			Metadata:       metadata,
		}

		err = ms.Write(r.Context(), cmd)
		if err != nil {
			log.Error("error writing to ms in posthandler: ", err)
		}
		json.NewEncoder(w).Encode(response)
	}
}
