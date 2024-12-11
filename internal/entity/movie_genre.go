package entity

type MovieGenre struct {
	ID      string `json:"id,omitempty" db:"id,omitempty"`
	MovieID string `json:"movie_id,omitempty" db:"movie_id,omitempty"`
	GenreID string `json:"genre_id,omitempty" db:"genre_id,omitempty"`
	DefaultCompleteTimestamp
}
