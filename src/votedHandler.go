package main

import (
	"context"

	gms "github.com/blackhatbrigade/gomessagestore"
	log "github.com/sirupsen/logrus"
)

type votedEventHandler struct {
	Repo      Repository
	Projector gms.Projector
}

func NewEventHandler(repo Repository, proj gms.Projector) votedEventHandler {
	return votedEventHandler{
		Repo:      repo,
		Projector: proj,
	}
}

func (h votedEventHandler) Type() string {
	return "Voted"
}

func (h votedEventHandler) Process(ctx context.Context, msg gms.Message) error {
	//process events
	// verify msg is an event
	evt := msg.(gms.Event)

	//make sure I haven't handled
	// do projection
	vote, err := h.Projector.Run(
		ctx,
		"Vote",
		evt.EntityID,
	)
	if err != nil {
		log.Error("Error while handling event: ", err)
	}

	newVoteState := vote.(voteState)
	//make the vote state something we can store
	log.Info("newVoteState from votedHandler: ", newVoteState)

	for k, v := range newVoteState.Candidate {
		rm := readmodel{
			ID:             newVoteState.ID,
			CandidateName:  k,
			CandidateTotal: v,
		}
		// don't double count votes
		err = h.Repo.store(ctx, rm)
		if err != nil {
			log.Error("Error while storing model: ", err)
		}
		log.Info("vote stored")
	}

	return nil
}
