package entities

import (
	"gourmand.golang.arango/src/enums"
)


type Authentication struct {
	Id       uint   `json:"_id,omitempty"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     enums.Role `json:"role"`
}