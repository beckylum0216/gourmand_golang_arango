package entities

type Author struct {
	Id           string   `json:"_Id,omitempty"`
	Source       string `json:"source"`
	CreditString string `json:"credit_string"`
}
