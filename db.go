package main

import (
	"log"
	"context"
	"fmt"
	"os"
	"time"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	UserID int `json:"user_id"`
	Email int `json:"email"`
}

type List struct {
	ListID int `json:"list_id"`
	Title string `json:"title"`
	TimeCreated time.Time `json:"time_created"`
	TimeModified time.Time `json:"time_modified"`
	UserID int `json:"user_id"`
}

type Entry struct {
	EntryID int `json:"entry_id"`
	Name string `json:"name"`
	Score int `json:"score"`
	ListID int `json:"list_id"`
}

var pool *pgxpool.Pool

func initializeDB() {
	fmt.Println("initializing db")
	//urlExample := "postgres://username:password@localhost:5432/databaseName"
	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error parsing database URL: %v", err)
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	fmt.Println("printing pool")
	fmt.Println(pool)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to PostgreSQL database via pgx: %v", err)
	}

	fmt.Println("Successful database connectoin")
}

