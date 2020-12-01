package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq" //Init postqresQL driver
	log "github.com/sirupsen/logrus"
)

func main() {

	// ==========================================
	// connect to databases and handle migrations
	// ==========================================

	// connect to the dbs
	//rmdb, _ := connectToDB()

	//log.Info("running migrations")
	//if migErr := handleMigrations(rmdb); migErr != nil {
	//	log.Error("migration error: ", migErr)
	//}

	//repo := NewRepository(rmdb)

	// ========================================
	// start server with middleware
	// ========================================

	server := &http.Server{
		Handler: r,
		Addr:    "localhost:8080", // env var
	}
	log.Info("server up and running at ", server.Addr)
	log.Fatal(server.ListenAndServe())
	// listen for error or shutdowns here
}
