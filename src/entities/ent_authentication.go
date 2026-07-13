package entities


type Authentication struct {
	Id       string `json:"_id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}