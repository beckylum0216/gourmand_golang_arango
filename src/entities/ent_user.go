package entities

import "gourmand.golang.arango/src/enums"

type User struct {
	Id       string     `json:"_id,omitempty"`
	Key      string     `json:"_key,omitempty"`
	Username string     `json:"username"`
	Active   bool       `json:"active"`
	Role     enums.Role `json:"role"`
}
