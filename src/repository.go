package main

import (
	"context"
	"database/sql"

	log "github.com/sirupsen/logrus"
)

//NewRepository connects to the sql db
func NewRepository(db *sql.DB) Repository {
	return Repository{
		db: db,
	}
}

// project from event stream
type Repository struct {
	db *sql.DB
}

func (r Repository) store(ctx context.Context, vote readmodel) error {
	query := `INSERT INTO candidate (id, name, voteTotal) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE id = ?, name = ?, voteTotal = ?`
	log.Info("vote in repo: ", vote)
	_, err := r.db.ExecContext(
		ctx,
		query,
		vote.ID,
		vote.CandidateName,
		vote.CandidateTotal,
		vote.ID, // start of upsert
		vote.CandidateName,
		vote.CandidateTotal,
	)
	if err != nil {
		log.Error("Error while storing model")
		return err
	}

	return nil
}
