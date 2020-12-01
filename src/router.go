package main

import (
	// import a kafka struct defined by it's interface
	// import a postgres struct defined by an interface
	"github.com/rochinworks/go-home/kafka"
	"github.com/rochinworks/go-home/pg"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func httpRouter(pg pg.Controller, kafka kafka.Controller) (chi.Router, error) {
	// chi router is easy to use and lightweight
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//=========//
	// WEATHER //
	//=========//
	r.Post("/weather", WeatherHandler(pg, kafka))
	r.Get("/weather/status", WeatherHandler(pg, kafka))

	//=========//
	// GARDEN  //
	//=========//
	r.Post("/garden", GardenHandler(pg, kafka))
	r.Get("/garden/status", GardenHandler(pg, kafka))

	//=========//
	// HOME 	 //
	//=========//
	r.Post("/home", HomeHandler(pg, kafka))
	r.Get("/home/status", HomeHandler(pg, kafka))

	//=========//
	// STORAGE //
	//=========//
	r.Post("/storage", StorageHandler(pg, kafka))

	r.Get("/storage/video/{id}", StorageHandler(pg, kafka))
	r.Get("/storage/docs/{id}", StorageHandler(pg, kafka))
	r.Get("/storage/photos/{id}", StorageHandler(pg, kafka))
	r.Get("/storage/music/{id}", StorageHandler(pg, kafka))
	r.Get("/storage/movies/{id}", StorageHandler(pg, kafka))

	r.Post("/storage/video", StorageHandler(pg, kafka))
	r.Post("/storage/docs", StorageHandler(pg, kafka))
	r.Post("/storage/photos", StorageHandler(pg, kafka))
	r.Post("/storage/music", StorageHandler(pg, kafka))
	r.Post("/storage/movies", StorageHandler(pg, kafka))

	return r, nil
}
