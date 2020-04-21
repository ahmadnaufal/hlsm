package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

func (a App) prepareRequest(method, url string, data url.Values) *http.Request {
	var formBody io.Reader = nil
	if data != nil {
		formBody = strings.NewReader(data.Encode())
	}

	req, _ := http.NewRequest(method, url, formBody)
	req.Header.Add("Cookie", a.Config.Cookie)
	req.Header.Add("Referer", url)

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	}

	return req
}

func buildAddressPayload(cfg AppConfig) url.Values {
	address := url.Values{}

	// default params
	address.Set("agree", "1")
	address.Set("provinsi", "")
	address.Set("kota", "")
	address.Set("x", "90")
	address.Set("y", "13")

	// params from appconfig
	address.Set("receive_at", strconv.Itoa(cfg.CheckoutOptions.ReceiveAt))
	address.Set("paymethod", strconv.Itoa(cfg.CheckoutOptions.PaymentMethod)) // 4 = jeketi points
	address.Set("fullname", cfg.User.Name)
	address.Set("address1", cfg.User.Address)
	address.Set("zipcode", cfg.User.Zipcode)
	address.Set("country", cfg.User.Country)
	address.Set("phone", cfg.User.Phone)

	return address
}

func main() {
	client := &http.Client{}
	app := loadApp(client, "payload.json")

	urls := prepareURLs(app.Config.HID, app.Config.PID)

	frm := url.Values{}
	frm.Set("box_7603", "1")
	frm.Set("x", "144")
	frm.Set("y", "12")

	req := app.prepareRequest("POST", urls.URLForm, frm)
	formResp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer formResp.Body.Close()

	formBody, _ := ioutil.ReadAll(formResp.Body)
	fmt.Println(string(formBody))

	address := buildAddressPayload(app.Config)
	req = app.prepareRequest("POST", urls.URLAddress, address)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	addressBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(addressBody))

	creq := app.prepareRequest("GET", urls.URLComply, nil)
	compResp, err := client.Do(creq)
	if err != nil {
		log.Println(err)
		return
	}
	defer compResp.Body.Close()

	compBody, _ := ioutil.ReadAll(compResp.Body)
	fmt.Println(string(compBody))
}
