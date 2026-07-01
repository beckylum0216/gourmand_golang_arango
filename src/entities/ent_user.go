package entities

type User struct {
	Id			   		string `json:"_id,omitempty"`
	Authentication   	Authentication
	Active           	bool `json:"active"`
}
