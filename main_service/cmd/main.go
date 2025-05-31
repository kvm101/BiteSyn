package main

import (
	"context"
	"log"
	"net/http"
	"restaurant_reviews/database"
	"restaurant_reviews/routes"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	err := database.ConnectMongo(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Run database migrations
	err = database.RunMigrations(ctx)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	r := routes.SetupRoutes()

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
