package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Mailer Mail
}

const (
	webPort = "80"
)

func main() {
	app := Config{
		Mailer: newMail(),
	}

	log.Printf("Starting mail service on port %s", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panicf("Error starting mail service on port %s, error: %v", webPort, err)
		return
	}
	log.Println("Mail service is running on port ", webPort)

}

func newMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
		FromName:    os.Getenv("FROM_NAME"),
	}
	return m
}
