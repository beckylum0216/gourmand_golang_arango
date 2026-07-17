package entities

type PersonWithDetails struct {
	Person Person  `json:"person"`
	User   *User   `json:"user,omitempty"`
	Author *Author `json:"author,omitempty"`
}

