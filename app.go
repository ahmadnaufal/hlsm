package hlsm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type App struct {
	Client *http.Client
	Config AppConfig
	URLs
}

type AppConfig struct {
	HID     int                          `json:"hid"`
	PID     int                          `json:"pid"`
	Tickets []map[string]map[string]uint `json:"tickets"`
	Cookie  string                       `json:"cookie"`

	User struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		Province string `json:"province"`
		City     string `json:"city"`
		Zipcode  string `json:"zipcode"`
		Country  string `json:"country"`
		Phone    string `json:"phone"`
	} `json:"user"`

	CheckoutOptions struct {
		ReceiveAt     int `json:"receive_at"`
		PaymentMethod int `json:"payment_method"`
	}
}

// URLs is the structs to store all needed urls
type URLs struct {
	URLForm    string
	URLAddress string
	URLComply  string
}

const (
	urlFormTemplate    = "https://jkt48.com/handshake/form/hid/%d/pid/%d"
	urlAddressTemplate = "https://jkt48.com/handshake/address/hid/%d/pid/%d"
	urlComplyTemplate  = "https://jkt48.com/handshake/comp/hid/%d/pid/%d"
)

func prepareURLs(hid, pid int) URLs {
	return URLs{
		URLForm:    fmt.Sprintf(urlFormTemplate, hid, pid),
		URLAddress: fmt.Sprintf(urlAddressTemplate, hid, pid),
		URLComply:  fmt.Sprintf(urlComplyTemplate, hid, pid),
	}
}

func New(client *http.Client, pathToJSON string) App {
	jsonFile, err := os.Open(pathToJSON)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	var cfg AppConfig
	if err := json.NewDecoder(jsonFile).Decode(&cfg); err != nil {
		panic(err)
	}

	urls := prepareURLs(cfg.HID, cfg.PID)

	return App{
		Client: client,
		Config: cfg,
		URLs:   urls,
	}
}
