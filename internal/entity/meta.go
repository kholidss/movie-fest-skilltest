package entity

type MetaPagination struct {
	Page  int `json:"page" db:"page,omitempty"`
	Limit int `json:"limit" db:"limit,omitempty"`
}
