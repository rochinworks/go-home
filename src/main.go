package main

import (
	"context"
	"net/http"

	gms "github.com/blackhatbrigade/gomessagestore"
	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq" //Init postqresQL driver
	log "github.com/sirupsen/logrus"
)

func main() {

	// ===================================
	// connect to databases
	// ===================================

	// connect to the dbs
	rmdb, msDB := connectToDB()

	log.Info("running migrations")
	if migErr := handleMigrations(rmdb); migErr != nil {
		log.Error("migration error: ", migErr)
	}

	repo := NewRepository(rmdb)
	ms := gms.NewMessageStore(msDB, log.New())

	// initialize a panic watcher
	defer recover()

	// Create a projector
	projector, err := ms.CreateProjector(
		//set default state of the projections
		gms.DefaultState(voteState{}), //state to use for the projection
		//reducers
		gms.WithReducer(&voteReducer{}),
	)
	if err != nil {
		log.WithError(err).Error("An error occurred while creating a projector")
	}

	// create a new subscriber with the handlers
	subscriber, err := ms.CreateSubscriber(
		"1",
		[]gms.MessageHandler{
			//handlers are connected here and projector is injected in
			NewHandler(projector, ms),
		},
		gms.SubscribeToCommandCategory("Vote"),
		gms.PollTime(100),
	)
	if err != nil {
		log.WithError(err).Error("An error occurred while creating a subscriber")
	}

	// event subscriber goes here
	votedsubscriber, err := ms.CreateSubscriber(
		"2",
		[]gms.MessageHandler{
			//handlers are connected here and projector is injected in
			NewEventHandler(repo, projector),
		},
		gms.SubscribeToCategory("Vote"),
		gms.PollTime(100),
	)
	if err != nil {
		log.WithError(err).Error("An error occurred while creating a subscriber")
	}

	// start poller loop inside a goroutine
	ctx := context.Background()
	go subscriber.Start(ctx)
	log.Info("subscriber started for Command stream")
	go votedsubscriber.Start(ctx)
	log.Info("subscriber started for Event stream")

	// ========================================
	// start server with recovery middleware
	// ========================================

	// chi router is easy to use and lightweight
	r := chi.NewRouter()

	r.Post("/vote", PostHandler(repo, ms))

	server := &http.Server{
		Handler: r,
		Addr:    "localhost:8080",
	}
	log.Info("server up and running at ", server.Addr)
	log.Fatal(server.ListenAndServe())
}
