package structs

type Groups struct {
	Id          string
	Creator     string
	Title       string
	Description string
	Date        string
}

type NewGroup struct {
	Creator     string
	Title       string
	Description string
}

var GroupID string