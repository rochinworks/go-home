package main

import (
	"encoding/json"
	//"fmt"

	gms "github.com/blackhatbrigade/gomessagestore"
	log "github.com/sirupsen/logrus"
)

type voteReducer struct{}

func (r *voteReducer) Type() string {
	return "Voted"
}

func (r *voteReducer) Reduce(msg gms.Message, previousState interface{}) interface{} {
	// set current state from previous state
	state := previousState.(voteState)
	// verify msg is an event
	event := msg.(gms.Event)

	// unpack metadata
	var eventMetadata metadata
	err := json.Unmarshal(event.Metadata, &eventMetadata)
	if err != nil {
		log.Error("Err while unmarshalling metadata in reducer")
	}

	// unpack data
	var eventData voteEvt
	err = json.Unmarshal(event.Data, &eventData)
	if err != nil {
		log.Error("Err while unmarshalling data in reducer")
	}

	candidateVote := make(map[string]int)
	candidateVote[eventData.Candidate] = 1

	// copy over new state
	state.ID = event.EntityID
	state.Candidate = candidateVote

	return state
}
