package entities

type PersonWithDetails struct {
	Person Person  `json:"person"`
	User   *User   `json:"user"`
	Author *Author `json:"author"`
}

type PeopleWithDetails struct {
	Person Person  `json:"person"`
	User   *User   `json:"user,omitempty"`
	Author *Author `json:"author,omitempty"`
}
