package presentation

type CreatedBy struct {
	ID       string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Entity   string `json:"entity,omitempty"`
}
