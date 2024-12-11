package entity

type User struct {
	ID       string `json:"id,omitempty" db:"id,omitempty"`
	FullName string `json:"full_name,omitempty" db:"full_name,omitempty"`
	Email    string `json:"email,omitempty" db:"email,omitempty"`
	Salary   int    `json:"salary,omitempty" db:"salary,omitempty"`
	Password string `json:"password,omitempty" db:"password,omitempty"`
	Entity   string `json:"entity,omitempty" db:"entity,omitempty"`
	DefaultCompleteTimestamp
}
