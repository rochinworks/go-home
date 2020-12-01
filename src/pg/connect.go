package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func connectToDB() (*sql.DB, *sql.DB) {
	msChan := make(chan *sql.DB, 1)
	mysqlChan := make(chan *sql.DB, 1)

	go func() {

		sqlHost := os.Getenv("RMDB_HOST")
		sqlDSN := fmt.Sprintf("mysql:example@tcp(%s:3306)/names", sqlHost)
		rmdb, err := sql.Open("mysql", sqlDSN)
		if err != nil {
			log.Fatal(fmt.Errorf("error connecting to mysql rmdb %+v", err))
		}

		log.Info("Pinging the rmdb")
		for {
			if rmdbErr := rmdb.Ping(); rmdbErr != nil {
				log.Errorf("an error occurred connecting to the rmdb trying again in 20 seconds: %v\n", rmdbErr)
				time.Sleep(time.Second * 20)
			} else {
				log.Info("connected to mysql db")
				break
			}
		}

		mysqlChan <- rmdb
	}()

	go func() {
		msHost := os.Getenv("MS_HOST")
		// this string normally comes from the config (environment var)
		pgDSN := fmt.Sprintf("user=message_store dbname=message_store sslmode=disable host=%s", msHost)
		// connect to the postgres DB
		db, err := sql.Open("postgres", pgDSN)
		if err != nil {
			log.Fatal(fmt.Errorf("error connecting to ms %+v", err))
		}

		// Ping the db to make sure we connected properly
		log.Info("Pinging the ms")
		for {
			if err := db.Ping(); err != nil {
				log.Errorf("an error occurred connecting to the messagestore trying again in 30 seconds: %v\n", err)
				time.Sleep(time.Second * 30)
			} else {
				log.Info("connected to ms db")
				break
			}
		}

		msChan <- db
	}()

	// wait for both DBs to be setup
	var msDB *sql.DB
	var mysqlDB *sql.DB
	done := false
	for !done {
		select {
		case db := <-msChan:
			if msDB == nil {
				msDB = db
				if mysqlDB != nil {
					done = true
				}
			}
		case db := <-mysqlChan:
			if mysqlDB == nil {
				mysqlDB = db
				if msDB != nil {
					done = true
				}
			}
		}
	}

	// instantiate a new message store using the connector library
	return mysqlDB, msDB
}
