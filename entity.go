package hlsm

type Ticket struct {
	Name     string `json:"name"`
	FormName string `json:"form_name"`
	Session  string `json:"session"`
	Quantity uint   `json:"quantity"`
}
