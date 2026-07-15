package entities

type PersonWithDetails struct {
	Person Person  `json:"person"`
	User   *User   `json:"user"`
	Author *Author `json:"author"`
}
