package model

type Segment struct {
	Id      int    `json:"id"`
	Slug    string `json:"slug"`
	Percent *int   `json:"percent,omitempty"`
}
