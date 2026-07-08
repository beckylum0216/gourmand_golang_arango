package entities

type Author struct {
	Id           string `json:"_id,omitempty"`
	Key          string `json:"_key,omitempty"`
	Source       string `json:"source"`
	CreditString string `json:"credit_string"`
}
