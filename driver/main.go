package main

import (
	"log"
	"net/http"

	"github.com/ahmadnaufal/hlsm"
)

func main() {
	client := &http.Client{}
	app := hlsm.New(client, "payload.json")

	tickets, err := app.GetTicketDetails()
	if err != nil {
		log.Println(err)
		return
	}

	err = app.PostTickets(tickets)
	if err != nil {
		log.Println(err)
		return
	}
}
