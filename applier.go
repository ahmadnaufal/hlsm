package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type App struct {
	Client *http.Client
	Config AppConfig
}

type AppConfig struct {
	HID     int                       `json:"hid"`
	PID     int                       `json:"pid"`
	Tickets map[string]map[string]int `json:"tickets"`
	Cookie  string                    `json:"cookie"`

	User struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Zipcode string `json:"zipcode"`
		Country string `json:"country"`
		Phone   string `json:"phone"`
	} `json:"user"`

	CheckoutOptions struct {
		ReceiveAt     int `json:"receive_at"`
		PaymentMethod int `json:"payment_method"`
	}
}

func loadApp(client *http.Client, pathToJSON string) App {
	jsonFile, err := os.Open(pathToJSON)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	var cfg AppConfig
	if err := json.NewDecoder(jsonFile).Decode(&cfg); err != nil {
		panic(err)
	}

	return App{
		client,
		cfg,
	}
}
