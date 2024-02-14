package main

import (
	"database/sql"
	"fmt"
	"github.com/f1zm0n/authentication/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

const port = "80"

var counts int64

const sleepSeconds = 2

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Printf("starting authentication service on port %s", port)

	db := connectToDB()
	if db == nil {
		log.Panic("Can't connect to postgres via authentication service")
	}
	app := Config{
		DB:     db,
		Models: data.New(db),
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panicf("error starting authentication service on port %s, err: %v", port, err)
	}
	log.Printf("authentications service is running on port %s", port)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready ...")
			counts++
		} else {
			log.Println("connected to postgres")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Printf("Backing off for %d seconds", sleepSeconds)
		time.Sleep(sleepSeconds * time.Second)
		continue
	}
}
