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

	for _, ticket := range tickets {
		log.Println(ticket.String())
	}

	err = app.PostTickets(tickets)
	if err != nil {
		log.Println(err)
		return
	}

	err = app.PostAddress()
	if err != nil {
		log.Println(err)
		return
	}

	err = app.PostFinal()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Pembelian tiket berhasil.")
}
