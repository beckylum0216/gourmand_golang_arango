package entities

import "time"

type Person struct {
	Id        string    `json:"_id,omitempty"`
	Key       string    `json:"_key,omitempty"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
