package main

import (
	"fmt"
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

	fmt.Println(tickets)

	err = app.PostTickets(tickets)
	if err != nil {
		log.Println(err)
		return
	}
}
