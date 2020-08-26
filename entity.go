package hlsm

import "fmt"

type Ticket struct {
	Name     string `json:"name"`
	FormName string `json:"form_name"`
	Session  uint   `json:"session"`
	Quantity uint   `json:"quantity"`
}

func (t Ticket) String() string {
	return fmt.Sprintf("Member: %s, Sesi: %d, Jumlah Tiket: %d", t.Name, t.Session, t.Quantity)
}
