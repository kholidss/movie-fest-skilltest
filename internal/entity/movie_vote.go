package entity

type MovieVote struct {
	ID      string `json:"id,omitempty" db:"id,omitempty"`
	UserID  string `json:"user_id,omitempty" db:"user_id,omitempty"`
	MovieID string `json:"movie_id,omitempty" db:"movie_id,omitempty"`
	DefaultCompleteTimestamp
}
