package entity

type Genre struct {
	ID         string `json:"id,omitempty" db:"id,omitempty"`
	Name       string `json:"name,omitempty" db:"name,omitempty"`
	Slug       string `json:"slug,omitempty" db:"slug,omitempty"`
	ViewNumber int    `json:"view_number,omitempty" db:"view_number,omitempty"`
	DefaultCompleteTimestamp
}
