package hlsm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

func (a App) GetTicketDetails() ([]Ticket, error) {
	req := a.prepareRequest("GET", a.URLForm, nil)
	formResp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer formResp.Body.Close()

	container, title, err := readParseHTMLBody(formResp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[GET-TicketDetails] On page: %s", title)

	// For each item found
	tickets := []Ticket{}
	for _, ticket := range a.Config.Tickets {
		for member := range ticket {
			for sesi := range ticket[member] {
				sesiNum, _ := strconv.ParseInt(sesi, 10, 32)
				sesiName := fmt.Sprintf("Sesi%d", sesiNum)
				sesiQty := ticket[member][sesi]

				rowQuery := `div:contains("` + sesiName + `") > a:contains("` + member + `")`
				checker := container.Find(rowQuery)
				if checker.Length() < 1 {
					err = fmt.Errorf("Error: Tiket untuk member dengan kata kunci \"%s\" sesi %s tidak ditemukan", member, sesi)
					return nil, err
				}

				memberName := checker.Text()
				formName, _ := checker.Parent().Parent().Find(`.formRight select`).Attr("name")

				tickets = append(tickets, Ticket{
					Name:     memberName,
					FormName: formName,
					Session:  uint(sesiNum),
					Quantity: sesiQty,
				})
			}
		}
	}

	return tickets, nil
}

func (a App) PostTickets(tickets []Ticket) error {
	form := url.Values{}
	for _, ticket := range tickets {
		form.Add(ticket.FormName, strconv.FormatUint(uint64(ticket.Quantity), 10))
	}
	form.Add("x", "144")
	form.Add("y", "13")

	req := a.prepareRequest("POST", a.URLForm, form)
	formResp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer formResp.Body.Close()

	container, title, err := readParseHTMLBody(formResp.Body)
	if err != nil {
		return err
	}

	log.Printf("[POST-Tickets] On page: %s", title)

	paymentMethodDropdown := container.Find("select#paymethod")
	if strings.TrimSpace(paymentMethodDropdown.Text()) == "" {
		return fmt.Errorf("Error: Kamu sudah tidak dapat membeli tiket untuk HS ID: %d", a.Config.HID)
	}

	return nil
}

func (a App) PostAddress() error {
	address := a.buildAddressPayload()

	req := a.prepareRequest("POST", a.URLAddress, address)
	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, title, err := readParseHTMLBody(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("[POST-Address] On page: %s", title)

	return nil
}

func (a App) PostFinal() error {
	req := a.prepareRequest("GET", a.URLComply, nil)
	compResp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	defer compResp.Body.Close()

	_, title, err := readParseHTMLBody(compResp.Body)
	if err != nil {
		return err
	}

	log.Printf("[POST-Comply] On page: %s", title)

	return nil
}

func (a App) buildAddressPayload() url.Values {
	address := url.Values{}

	// default params
	address.Set("agree", "1")
	address.Set("x", "90")
	address.Set("y", "13")

	// params from appconfig
	address.Set("receive_at", strconv.Itoa(a.Config.CheckoutOptions.ReceiveAt))
	address.Set("paymethod", strconv.Itoa(a.Config.CheckoutOptions.PaymentMethod)) // 4 = jeketi points
	address.Set("fullname", a.Config.User.Name)
	address.Set("address1", a.Config.User.Address)
	address.Set("zipcode", a.Config.User.Zipcode)
	address.Set("country", a.Config.User.Country)
	address.Set("phone", a.Config.User.Phone)
	address.Set("provinsi", a.Config.User.Province)
	address.Set("kota", a.Config.User.City)

	return address
}
