package main

import (
	"context"
	"encoding/json"

	gms "github.com/blackhatbrigade/gomessagestore"
	log "github.com/sirupsen/logrus"
)

// gms.Handler should have a Type() and Process() method to satisfy the interface
// additionally command handlers should have a projector, and ms dependencies, so that we can
// inject those at compile time
type voteCommandHandler struct {
	Projector gms.Projector
	MS        gms.MessageStore
}

func NewHandler(projector gms.Projector, ms gms.MessageStore) voteCommandHandler {
	return voteCommandHandler{
		Projector: projector,
		MS:        ms,
	}
}

// command handler Type() should be the command messageType
func (h voteCommandHandler) Type() string {
	return "VoteForCandidate"
}

func (h voteCommandHandler) Process(ctx context.Context, msg gms.Message) error {
	// make sure it's a command (type assertion) otherwise we panic
	cmd := msg.(gms.Command)

	// Unmarshal metadata
	var meta metadata
	err := json.Unmarshal(cmd.Metadata, &meta)
	if err != nil {
		log.Error("and error occurred while unmarshalling metadata: ", err)
		return err
	}
	log.Info("metadata unmarshalled")

	// project on event stream
	newVoteState, err := h.Projector.Run(
		ctx,
		"Vote",
		meta.UserID,
	)
	log.Info("newVoteState projection from command handler: ", newVoteState)
	if err != nil {
		return err
	}

	var voteEventData voteEvt
	err = json.Unmarshal(cmd.Data, &voteEventData)
	if err != nil {
		log.Error("and error occurred while unmarshalling data: ", err)
		return err
	}
	log.Info("data unmarshalled")

	checkedVoteState := newVoteState.(voteState)

	// check to see if this is the first message in the stream otherwise
	// don't count the vote
	if checkedVoteState.ID != gms.NilUUID {
		log.Debug("vote has already been counted")
		return nil
	}
	log.Info("Checked vote state: ", checkedVoteState)
	log.Info("vote from command: ", voteEventData.Candidate)

	evt := &gms.Event{
		ID:             gms.NewID(),
		EntityID:       meta.UserID,
		StreamCategory: "Vote",
		MessageType:    "Voted", // voted is the message type of the event being written
		Data:           cmd.Data,
		Metadata:       cmd.Metadata,
	}

	// write event
	err = h.MS.Write(ctx, evt)
	if err != nil {
		return err
	}
	log.Info("command handled")
	return nil

}
