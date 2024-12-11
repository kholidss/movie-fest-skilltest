package entity

type ActionHistory struct {
	ID             string `json:"id,omitempty" db:"id,omitempty"`
	Name           string `json:"name,omitempty" db:"name,omitempty"`
	IdentifierID   string `json:"identifier_id,omitempty" db:"identifier_id,omitempty"`
	IdentifierName string `json:"identifier_name,omitempty" db:"identifier_name,omitempty"`
	UserAgent      string `json:"user_agent,omitempty" db:"user_agent,omitempty"`
}
