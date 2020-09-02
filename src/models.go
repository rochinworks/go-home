package main

import (
	"github.com/blackhatbrigade/gomessagestore/uuid"
)

// repo structure
type readmodel struct {
	ID             uuid.UUID `db:"id"`
	CandidateName  string    `db:"name"`
	CandidateTotal int       `db:"voteTotal"`
}

// default state
type voteState struct {
	ID        uuid.UUID
	Candidate map[string]int `json:"candidate"`
}

// command
type voteCmd struct {
	Candidate string `json:"candidate"`
}

// event
type voteEvt struct {
	Candidate string `json:"candidate"`
}

// metadata
type metadata struct {
	UserID uuid.UUID `json:"userId"`
}
